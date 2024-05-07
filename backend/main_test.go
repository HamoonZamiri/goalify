package main_test

import (
	"bytes"
	"encoding/json"
  "goalify"
	"goalify/db"
	"goalify/entities"
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
	reqBody := handler.SignupRequest{Email: "user@mail.com", Password: "password"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, err := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
	var serverResponse responses.ServerResponse[entities.User]

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	json.Unmarshal(body, &serverResponse)
	assert.Equal(t, "user@mail.com", serverResponse.Data.Email)
}

func TestSignupEmailExists(t *testing.T) {
	reqBody := handler.SignupRequest{Email: "user@mail.com", Password: "password"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestLogin(t *testing.T) {
	reqBody := handler.LoginRequest{Email: "user@mail.com", Password: "password"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, _ := http.Post(BASE_URL+"/api/users/login", "application/json", bytes.NewReader(stringifiedBody))
	assert.Equal(t, res.StatusCode, http.StatusOK)
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
