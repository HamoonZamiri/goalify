package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goalify/db"
	"goalify/entities"
	gh "goalify/goals/handler"
	"goalify/responses"
	"goalify/users/handler"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const BASE_URL = "http://localhost:8080"

var (
	accessToken  string
	refreshToken string
	userId       string
)

func setup() {
	go main.Run()
	time.Sleep(2 * time.Second)
}

func TestHealth(t *testing.T) {
	res, err := http.Get(BASE_URL + "/health")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestSignup(t *testing.T) {
	reqBody := handler.SignupRequest{Email: "user@mail.com", Password: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, err := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
	var serverResponse responses.ServerResponse[entities.UserDTO]

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	json.Unmarshal(body, &serverResponse)
	assert.Equal(t, "user@mail.com", serverResponse.Data.Email)
	refreshToken = serverResponse.Data.RefreshToken.String()
	userId = serverResponse.Data.Id.String()
}

func TestSignupEmailExists(t *testing.T) {
	reqBody := handler.SignupRequest{Email: "user@mail.com", Password: "password"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestLogin(t *testing.T) {
	reqBody := handler.LoginRequest{Email: "user@mail.com", Password: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/login", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestLoginIncorrectPassword(t *testing.T) {
	reqBody := handler.LoginRequest{Email: "user@mail.com", Password: "password123!!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/login", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRefresh(t *testing.T) {
	prevToken := refreshToken
	reqBody := handler.RefreshRequest{UserId: userId, RefreshToken: prevToken}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	url := fmt.Sprintf("%s/api/users/refresh", BASE_URL)
	res, err := http.Post(url, "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	var serverResponse responses.ServerResponse[entities.UserDTO]

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &serverResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	accessToken = serverResponse.Data.AccessToken
	refreshToken = serverResponse.Data.RefreshToken.String()
	assert.NotEqual(t, prevToken, serverResponse.Data.RefreshToken.String())
}

func TestIncorrectRefresh(t *testing.T) {
	reqBody := handler.RefreshRequest{UserId: userId, RefreshToken: "invalid"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	url := fmt.Sprintf("%s/api/users/refresh", BASE_URL)
	res, err := http.Post(url, "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGoalCategoryCreate(t *testing.T) {
	reqBody := gh.NewGoalCategoryRequest("goal cat", 100)
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	req.Header.Set("Authorization", "Bearer "+refreshToken)
}

func TestMain(m *testing.M) {
	db, err := db.New("goalify")
	if err != nil {
		panic(err)
	}
	setup()
	code := m.Run()

	query := `DELETE from users;`
	db.MustExec(query)
	os.Exit(code)
}
