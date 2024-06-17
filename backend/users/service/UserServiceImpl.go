package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/svcerror"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userStore      stores.UserStore
	eventPublisher events.EventPublisher
}

func NewUserService(userStore stores.UserStore, ep events.EventPublisher) *UserServiceImpl {
	return &UserServiceImpl{userStore: userStore, eventPublisher: ep}
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

func (s *UserServiceImpl) SignUp(email, password string) (*entities.UserDTO, error) {
	_, err := s.userStore.GetUserByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("%w: user with email %s already exists", svcerror.ErrBadRequest, email)
	}

	if err != sql.ErrNoRows {
		slog.Error("service.SignUp: store.GetUserByEmail:", "err", err.Error())
		return nil, fmt.Errorf("%w: internal error signing up user", svcerror.ErrInternalServer)
	}

	cleanedEmail := strings.TrimSpace(email)
	cleanedEmail = strings.ToLower(cleanedEmail)
	if cleanedEmail == "" {
		return nil, fmt.Errorf("%w: email cannot be empty", svcerror.ErrBadRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("service.SignUp: bcrypt.GenerateFromPassword:", "err", err.Error())
		return nil, fmt.Errorf("%w: error hashing password", svcerror.ErrInternalServer)
	}

	user, err := s.userStore.CreateUser(cleanedEmail, string(hashedPassword))
	if err != nil {
		slog.Error("service.SignUp: store.CreateUser:", "err", err.Error())
		return nil, fmt.Errorf("%w: error creating user", svcerror.ErrInternalServer)
	}
	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Refresh: service.SignUp:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", svcerror.ErrInternalServer)
	}
	s.eventPublisher.Publish(events.NewEvent("user_created", user))
	return user.ToUserDTO(accessToken), nil
}

func (s *UserServiceImpl) Refresh(userId, refreshToken string) (*entities.UserDTO, error) {
	user, err := s.userStore.GetUserById(userId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: error finding user", svcerror.ErrNotFound)
	}

	if err != nil {
		slog.Error("service.Refresh: store.GetUserById:", "err", err.Error())
		return nil, fmt.Errorf("%w: error getting user", svcerror.ErrInternalServer)
	}

	if user.RefreshToken.String() != refreshToken {
		return nil, fmt.Errorf("%w: invalid refresh token", svcerror.ErrBadRequest)
	}

	if user.RefreshTokenExpiry.Before(time.Now()) {
		return nil, fmt.Errorf("%w: refresh token expired", svcerror.ErrBadRequest)
	}

	newRefreshToken := generateRefreshToken()
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), newRefreshToken.String())
	if err != nil {
		slog.Error("service.Refresh: store.UpdateRefreshToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error updating refresh token", svcerror.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Refresh: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", svcerror.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *UserServiceImpl) Login(email, password string) (*entities.UserDTO, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: user with email %s not found", svcerror.ErrNotFound, email)
	} else if err != nil {
		slog.Error("service.Login: store.GetUserByEmail:", "err", err.Error())
		return nil, fmt.Errorf("%w: error getting user", svcerror.ErrInternalServer)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("%w: invalid password", svcerror.ErrBadRequest)
	}

	if err != nil {
		slog.Error("service.Login: bcrypt.CompareHashAndPassword:", "err", err.Error())
		return nil, fmt.Errorf("%w: error comparing password", svcerror.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.Login: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", svcerror.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}

func (s *UserServiceImpl) DeleteUserById(id string) error {
	err := s.userStore.DeleteUserById(id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: user not found", svcerror.ErrNotFound)
	}
	if err != nil {
		slog.Error("service.DeleteUserById: store.DeleteUserById:", "err", err.Error())
		return fmt.Errorf("%w: error deleting user", svcerror.ErrInternalServer)
	}
	return nil
}

func (s *UserServiceImpl) UpdateUserById(id uuid.UUID, updates map[string]interface{}) (*entities.UserDTO, error) {
	user, err := s.userStore.UpdateUserById(id, updates)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: user not found", svcerror.ErrNotFound)
	}

	if err != nil {
		slog.Error("service.UpdateUserById: store.UpdateUserById", "err", err.Error())
		return nil, fmt.Errorf("%w: error updating user", svcerror.ErrInternalServer)
	}

	accessToken, err := generateJWTToken(user.Id)
	if err != nil {
		slog.Error("Users: service.UpdateUserById: service.generateJWTToken:", "err", err.Error())
		return nil, fmt.Errorf("%w: error generating access token", svcerror.ErrInternalServer)
	}
	return user.ToUserDTO(accessToken), nil
}
