package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	main "goalify"
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

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const BASE_URL = "http://localhost:8080"

var (
	dbx          *sqlx.DB
	accessToken  string
	refreshToken string
	userId       string
)

func MarshalServerResponse[T any](res *http.Response) (responses.ServerResponse[T], error) {
	var serverResponse responses.ServerResponse[T]

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return serverResponse, err
	}

	err = json.Unmarshal(body, &serverResponse)
	return serverResponse, err
}

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

	defer res.Body.Close()

	serverResponse, err := MarshalServerResponse[entities.UserDTO](res)
	assert.Nil(t, err)
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

	serverResponse, err := MarshalServerResponse[entities.UserDTO](res)
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

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	gc, err := MarshalServerResponse[entities.GoalCategory](res)

	assert.Nil(t, err)
	assert.Equal(t, "goal cat", gc.Data.Title)
	assert.Equal(t, 100, gc.Data.Xp_per_goal)
}

func TestGetGoalCategories(t *testing.T) {
	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := MarshalServerResponse[[]entities.GoalCategory](res)
	assert.Nil(t, err)

	t.Log(resBody)

	assert.Equal(t, 1, len(resBody.Data))
}

func TestGetGoalCategoryById(t *testing.T) {
	gc := createTestGoalCategory(t, uuid.MustParse(userId))

	url := fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := MarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)

	assert.Equal(t, gc.Id, resBody.Data.Id)
}

func printErrResponse(t *testing.T, res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var errRes responses.APIError
	err = json.Unmarshal(body, &errRes)
	if err != nil {
		return err
	}

	t.Logf("Error: %s", errRes.Message)
	return nil
}

func printSuccessResponse[T any](t *testing.T, res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var successRes responses.ServerResponse[T]
	err = json.Unmarshal(body, &successRes)
	if err != nil {
		return err
	}
	t.Logf("Response: %v", successRes)
	return nil
}

func createTestGoalCategory(t *testing.T, userId uuid.UUID) *entities.GoalCategory {
	query := `INSERT INTO goal_categories (title, xp_per_goal, user_id) 
  VALUES ($1, $2, $3) RETURNING *`
	var goalCategory entities.GoalCategory
	dbx.QueryRowx(query, "test goal", 50, userId).StructScan(&goalCategory)
	return &goalCategory
}

func TestMain(m *testing.M) {
	var err error
	dbx, err = db.New("goalify")
	if err != nil {
		panic(err)
	}

	setup()
	code := m.Run()

	query := `DELETE from goal_categories; DELETE from goals; DELETE from users;`
	dbx.MustExec(query)
	os.Exit(code)
}
