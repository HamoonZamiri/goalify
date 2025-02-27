package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/responses"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var subscribedEvents = []string{events.GOAL_UPDATED}

type UserService interface {
	SignUp(email, password string) (*entities.UserDTO, error)
	Login(email, password string) (*entities.UserDTO, error)
	Refresh(email, refreshToken string) (*entities.UserDTO, error)
	DeleteUserById(id string) error
	UpdateUserById(id uuid.UUID, updates map[string]interface{}) (*entities.UserDTO, error)

	GetLevelById(id int) (*entities.Level, error)
}

type userService struct {
	userStore      stores.UserStore
	eventPublisher events.EventPublisher
}

func NewUserService(userStore stores.UserStore, ep events.EventPublisher) UserService {
	us := &userService{userStore: userStore, eventPublisher: ep}

	for _, event := range subscribedEvents {
		ep.Subscribe(event, us)
	}
	return us
}

func generateJWTToken(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	return claims["userId"].(string), nil
}

func generateRefreshToken() uuid.UUID {
	return uuid.New()
}

func (s *userService) SignUp(email, password string) (*entities.UserDTO, error) {
	_, err := s.userStore.GetUserByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("%w: user with email %s already exists", responses.ErrBadRequest, email)
	}

	if err != sql.ErrNoRows {
		slog.Error("service.SignUp: store.GetUserByEmail:", "err", err.Error())
		return nil, fmt.Errorf("%w: internal error signing up user", responses.ErrInternalServer)
	}

	cleanedEmail := strings.TrimSpace(email)
	if cleanedEmail == "" {
		return nil, fmt.Errorf("%w: email cannot be empty", responses.ErrBadRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("service.SignUp: bcrypt.GenerateFromPassword:", "err", err.Error())
		return nil, fmt.Errorf("%w: error hashing password", responses.ErrInternalServer)
	}

	user, err := s.userStore.CreateUser(cleanedEmail, string(hashedPassword))
	if err != nil {
		slog.Error("service.SignUp: store.CreateUser:", "err", err.Error())
		return nil, fmt.Errorf("%w: error creating user", responses.ErrInternalServer)
	}
	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Refresh: service.SignUp:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", responses.ErrInternalServer)
	}
	s.eventPublisher.Publish(events.NewEventWithUserId(events.USER_CREATED, user, user.Id.String()))
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) Refresh(userId, refreshToken string) (*entities.UserDTO, error) {
	user, err := s.userStore.GetUserById(userId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: error finding user", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error("service.Refresh: store.GetUserById:", "err", err.Error())
		return nil, fmt.Errorf("%w: error getting user", responses.ErrInternalServer)
	}

	if user.RefreshToken.String() != refreshToken {
		return nil, fmt.Errorf("%w: invalid refresh token", responses.ErrBadRequest)
	}

	if user.RefreshTokenExpiry.Before(time.Now()) {
		return nil, fmt.Errorf("%w: refresh token expired", responses.ErrBadRequest)
	}

	newRefreshToken := generateRefreshToken()
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), newRefreshToken.String())
	if err != nil {
		slog.Error("service.Refresh: store.UpdateRefreshToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error updating refresh token", responses.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Refresh: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", responses.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) Login(email, password string) (*entities.UserDTO, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: user with email %s not found", responses.ErrNotFound, email)
	} else if err != nil {
		slog.Error("service.Login: store.GetUserByEmail:", "err", err.Error())
		return nil, fmt.Errorf("%w: error getting user", responses.ErrInternalServer)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("%w: invalid password", responses.ErrBadRequest)
	}

	if err != nil {
		slog.Error("service.Login: bcrypt.CompareHashAndPassword:", "err", err.Error())
		return nil, fmt.Errorf("%w: error comparing password", responses.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Login: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", responses.ErrInternalServer)
	}

	// refresh token now
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), generateRefreshToken().String())
	if err != nil {
		slog.Error("service.Login: store.UpdateRefreshToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error updating refresh token", responses.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) DeleteUserById(id string) error {
	err := s.userStore.DeleteUserById(id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: user not found", responses.ErrNotFound)
	}
	if err != nil {
		slog.Error("service.DeleteUserById: store.DeleteUserById:", "err", err.Error())
		return fmt.Errorf("%w: error deleting user", responses.ErrInternalServer)
	}
	return nil
}

func (s *userService) UpdateUserById(id uuid.UUID, updates map[string]interface{}) (*entities.UserDTO, error) {
	user, err := s.userStore.UpdateUserById(id, updates)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: user not found", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error("service.UpdateUserById: store.UpdateUserById", "err", err.Error())
		return nil, fmt.Errorf("%w: error updating user", responses.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.UpdateUserById: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", responses.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) GetLevelById(id int) (*entities.Level, error) {
	level, err := s.userStore.GetLevelById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: could not get level information", responses.ErrBadRequest)
	}
	return level, nil
}
