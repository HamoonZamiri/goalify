package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"goalify/internal/entities"
	"goalify/internal/users/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	reqBody := handler.SignupRequest{Email: email, Password: "password123!", ConfirmPassword: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, err := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	defer res.Body.Close()

	serverResponse, err := unmarshalResponse[entities.UserDTO](res)
	assert.Nil(t, err)
	assert.Equal(t, email, serverResponse.Email)
}

func TestSignupEmailExists(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	createUser(email, "password123!")
	reqBody := handler.SignupRequest{Email: email, Password: "password123!", ConfirmPassword: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestLogin(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	createUser(email, "password123!")
	reqBody := handler.LoginRequest{Email: email, Password: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/login", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Nil(t, err)

	user, err := unmarshalResponse[entities.UserDTO](res)
	assert.Nil(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, email, user.Email)
}

func TestLoginIncorrectPassword(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	createUser(email, "password123!")
	reqBody := handler.LoginRequest{Email: email, Password: "password123!!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/login", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestRefresh(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	user := createUser(email, "password123!")
	prevToken := user.RefreshToken
	reqBody := map[string]any{"user_id": user.Id, "refresh_token": prevToken}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	url := fmt.Sprintf("%s/api/users/refresh", BASE_URL)
	res, err := http.Post(url, "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	serverResponse, err := unmarshalResponse[entities.UserDTO](res)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.NotEqual(t, prevToken, serverResponse.RefreshToken.String())
}

func TestIncorrectRefresh(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	reqBody := map[string]any{"user_id": userDto.Id, "refresh": "incorrect"}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/users/refresh", BASE_URL), reqBody, userDto.AccessToken)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}

func TestUpdateUserById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	reqBody := map[string]any{"xp": 100, "cash_available": 100}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/users", BASE_URL), reqBody, userDto.AccessToken)
	require.Nil(t, err)

	resBody, err := unmarshalResponse[entities.UserDTO](res)
	assert.Nil(t, err)
	assert.Equal(t, 100, resBody.Xp)
	assert.Equal(t, 100, resBody.CashAvailable)
}

func TestIncorrectUpdateUserById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	reqBody := map[string]any{"xp": "incorrect", "cash_available": "incorrect"}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/users", BASE_URL), reqBody, userDto.AccessToken)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetLevelById(t *testing.T) {
	t.Parallel()
	user := createUser(t.Name()+"@mail.com", "Password123!")
	for level := 1; level <= 100; level++ {
		url := fmt.Sprintf("%s/api/levels/%d", BASE_URL, level)
		res, err := buildAndSendRequest("GET", url, nil, user.AccessToken)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}
