package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/internal/config"
	"goalify/internal/entities"
	"goalify/internal/users/stores"
	"goalify/internal/events"
	"goalify/internal/responses"
	"log/slog"
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
	VerifyToken(tokenString string) (string, error)

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

	tokenString, err := token.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) VerifyToken(tokenString string) (string, error) {
	errResponse := fmt.Errorf("%w: could not authenticate request", responses.ErrUnauthorized)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})
	if err != nil {
		slog.Error("service.VerifyToken: jwt.Parse:", "err", err.Error())
		return "", errResponse
	}

	if !token.Valid {
		return "", errResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errResponse
	}
	userId := claims["userId"]
	castedUserId, ok := userId.(string)
	if !ok {
		return "", errResponse
	}

	_, err = s.userStore.GetUserById(castedUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errResponse
	}
	if err != nil {
		slog.Error("service.VerifyToken: store.GetUserById:", "err", err.Error())
		return "", responses.ErrInternalServer
	}
	return castedUserId, nil
}

func generateRefreshToken() uuid.UUID {
	return uuid.New()
}

func (s *userService) SignUp(email, password string) (*entities.UserDTO, error) {
	_, err := s.userStore.GetUserByEmail(email)

	badReqErr := fmt.Errorf("%w: invalid signup request", responses.ErrBadRequest)
	if err == nil {
		return nil, badReqErr
	}

	if !errors.Is(err, sql.ErrNoRows) {
		slog.Error("service.SignUp: store.GetUserByEmail:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("service.SignUp: bcrypt.GenerateFromPassword:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	user, err := s.userStore.CreateUser(email, string(hashedPassword))
	if err != nil {
		slog.Error("service.SignUp: store.CreateUser:", "err", err.Error())
		return nil, fmt.Errorf("%w: error creating user", responses.ErrInternalServer)
	}
	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("service.SignUp:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}
	s.eventPublisher.Publish(events.NewEventWithUserId(events.USER_CREATED, user, user.Id.String()))
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) Refresh(userId, refreshToken string) (*entities.UserDTO, error) {
	errResponse := fmt.Errorf("%w: invalid user id or refresh token", responses.ErrBadRequest)
	user, err := s.userStore.GetUserById(userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errResponse
	}

	if err != nil {
		slog.Error("service.Refresh: store.GetUserById:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	if user.RefreshToken.String() != refreshToken {
		return nil, errResponse
	}

	if user.RefreshTokenExpiry.Before(time.Now()) {
		return nil, errResponse
	}

	newRefreshToken := generateRefreshToken()
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), newRefreshToken.String())
	if err != nil {
		slog.Error("service.Refresh: store.UpdateRefreshToken:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Refresh: service.generateJWTToken:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) Login(email, password string) (*entities.UserDTO, error) {
	errResponse := fmt.Errorf("%w: invalid email or password", responses.ErrBadRequest)
	user, err := s.userStore.GetUserByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errResponse
	}
	if err != nil {
		slog.Error("service.Login: store.GetUserByEmail:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, errResponse
	}

	if err != nil {
		slog.Error("service.Login: bcrypt.CompareHashAndPassword:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Login: service.generateJWTToken:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	// refresh token now
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), generateRefreshToken().String())
	if err != nil {
		slog.Error("service.Login: store.UpdateRefreshToken:", "err", err.Error())
		return nil, responses.ErrInternalServer
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *userService) DeleteUserById(id string) error {
	err := s.userStore.DeleteUserById(id)
	if errors.Is(err, sql.ErrNoRows) {
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
	// should not reach this point if the user does not exist
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("service.UpdateUserById: store.UpdateUserById", "err", err.Error())
		return nil, responses.ErrInternalServer
	}

	if err != nil {
		slog.Error("service.UpdateUserById: store.UpdateUserById", "err", err.Error())
		return nil, responses.ErrInternalServer
	}
	return user.ToUserDTO(""), nil
}

func (s *userService) GetLevelById(id int) (*entities.Level, error) {
	level, err := s.userStore.GetLevelById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: could not get level information", responses.ErrBadRequest)
	}
	return level, nil
}
