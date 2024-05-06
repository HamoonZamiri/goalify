package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/users/stores"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userStore stores.UserStore
}

func NewUserService(userStore stores.UserStore) *UserServiceImpl {
	return &UserServiceImpl{userStore: userStore}
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
		return os.Getenv("JWT_SECRET"), nil
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

func userToUserDTO(user entities.User) entities.UserDTO {
	token, err := generateJWTToken(user.Id)
	if err != nil {
		// handle error
		panic(err)
	}

	return entities.UserDTO{
		Email:              user.Email,
		AccessToken:        token,
		Xp:                 user.Xp,
		LevelId:            user.LevelId,
		CashAvailable:      user.CashAvailable,
		Id:                 user.Id,
		RefreshToken:       user.RefreshToken,
		RefreshTokenExpiry: user.RefreshTokenExpiry,
	}
}

func (s *UserServiceImpl) SignUp(email, password string) (entities.UserDTO, error) {
	_, err := s.userStore.GetUserByEmail(email)
	if err == nil {
		return entities.UserDTO{}, fmt.Errorf("user with email %s already exists", email)
	}

	if err != sql.ErrNoRows {
		return entities.UserDTO{}, err
	}

	cleanedEmail := strings.TrimSpace(email)
	cleanedEmail = strings.ToLower(cleanedEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserDTO{}, err
	}

	user, err := s.userStore.CreateUser(cleanedEmail, string(hashedPassword))
	if err != nil {
		return entities.UserDTO{}, err
	}
	return userDTOReturnVal(user, err)
}

func (s *UserServiceImpl) Refresh(userId, refreshToken string) (entities.UserDTO, error) {
	user, err := s.userStore.GetUserById(userId)
	if err != nil {
		return entities.UserDTO{}, err
	}
	if user.RefreshToken.String() != refreshToken {
		return entities.UserDTO{}, errors.New("invalid refresh token")
	}

	if user.RefreshTokenExpiry.Before(time.Now()) {
		return entities.UserDTO{}, errors.New("refresh token expired")
	}

	newRefreshToken := generateRefreshToken()
	user, err = s.userStore.UpdateRefreshToken(user.Id.String(), newRefreshToken.String())
	if err != nil {
		return entities.UserDTO{}, err
	}

	return userDTOReturnVal(user, err)
}

func userDTOReturnVal(user entities.User, err error) (entities.UserDTO, error) {
	return userToUserDTO(user), nil
}

func (s *UserServiceImpl) Login(email, password string) (entities.UserDTO, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err != nil {
		return entities.UserDTO{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return entities.UserDTO{}, errors.New("invalid password")
	}

	return userDTOReturnVal(user, err)
}
