package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	main "goalify"
	"goalify/config"
	"goalify/db"
	"goalify/entities"
	"goalify/users/handler"
	"goalify/utils/options"
	"goalify/utils/responses"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const BASE_URL = "http://localhost:8080"

var (
	dbx           *sqlx.DB
	accessToken   string
	refreshToken  string
	userId        string
	configService *config.ConfigService
)

func setup() {
	go main.Run()
	time.Sleep(250 * time.Millisecond)
}

func TestMain(m *testing.M) {
	var err error
	configService = config.NewConfigService(options.None[string]())
	configService.SetEnv("ENV", "test")

	dbx, err = db.New(configService.MustGetEnv("TEST_DB_NAME"),
		configService.MustGetEnv("DB_USER"), configService.MustGetEnv("DB_PASSWORD"))
	if err != nil {
		panic(err)
	}
	query := `DELETE from goals; DELETE FROM goal_categories; DELETE from users;`
	dbx.MustExec(query)

	setup()
	code := m.Run()

	dbx.MustExec(query)
	os.Exit(code)
}

func TestHealth(t *testing.T) {
	res, err := http.Get(BASE_URL + "/health")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

/* Users Domain Tests
* Testing Resource: /api/users
 */

func TestSignup(t *testing.T) {
	reqBody := handler.SignupRequest{Email: "user@mail.com", Password: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, err := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()

	serverResponse, err := UnmarshalServerResponse[entities.UserDTO](res)
	fmt.Println(io.ReadAll(res.Body))
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
	reqBody := map[string]any{"user_id": userId, "refresh_token": refreshToken}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	url := fmt.Sprintf("%s/api/users/refresh", BASE_URL)
	res, err := http.Post(url, "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	serverResponse, err := UnmarshalServerResponse[entities.UserDTO](res)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	accessToken = serverResponse.Data.AccessToken
	refreshToken = serverResponse.Data.RefreshToken.String()
	assert.NotEqual(t, prevToken, serverResponse.Data.RefreshToken.String())
}

func TestIncorrectRefresh(t *testing.T) {
	reqBody := map[string]any{"user_id": userId, "refresh": "incorrect"}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/users/refresh", BASE_URL), reqBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateUserById(t *testing.T) {
	reqBody := map[string]any{"xp": 100, "cash_available": 100}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/users", BASE_URL), reqBody)
	require.Nil(t, err)

	resBody, err := UnmarshalServerResponse[entities.UserDTO](res)
	assert.Nil(t, err)
	assert.Equal(t, 100, resBody.Data.Xp)
	assert.Equal(t, 100, resBody.Data.CashAvailable)
}

func TestIncorrectUpdateUserById(t *testing.T) {
	reqBody := map[string]any{"xp": "incorrect", "cash_available": "incorrect"}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/users", BASE_URL), reqBody)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

/* Goal Domain Tests
* Testing Resource: /api/goals
 */

func TestGoalCategoryCreate(t *testing.T) {
	reqBody := map[string]any{"title": "goal cat", "xp_per_goal": 100}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/goals/categories", BASE_URL), reqBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	gc, err := UnmarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, "goal cat", gc.Data.Title)
	assert.Equal(t, 100, gc.Data.Xp_per_goal)
	assert.NotNil(t, gc.Data.Goals)
	assert.Equal(t, userId, gc.Data.UserId.String())
}

func TestGoalCategoryCreateInvalidFields(t *testing.T) {
	reqBody := map[string]any{"title": "goal cat", "xp_per_goal": -10}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/goals/categories", BASE_URL), reqBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetGoalCategories(t *testing.T) {
	res, err := buildAndSendRequest("GET", fmt.Sprintf("%s/api/goals/categories", BASE_URL), nil)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := UnmarshalServerResponse[[]entities.GoalCategory](res)
	assert.Nil(t, err)

	// previously was asserting this value as 1, but since we have a user_created
	// event that creates a goal category, we should expect 2
	assert.Equal(t, 2, len(resBody.Data))
}

func TestGetGoalCategoryById(t *testing.T) {
	gc := createTestGoalCategory("create goal category", uuid.MustParse(userId))

	res, err := buildAndSendRequest("GET", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := UnmarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, gc.Id, resBody.Data.Id)
}

func TestUpdateGoalCategoryById(t *testing.T) {
	gc := createTestGoalCategory("update goal category", uuid.MustParse(userId))

	reqBody := map[string]any{"title": "updated title", "xp_per_goal": 69}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), reqBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := UnmarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, "updated title", resBody.Data.Title)
	assert.Equal(t, 69, resBody.Data.Xp_per_goal)
}

func TestUpdateGoalCategoryByIdInvalidFields(t *testing.T) {
	gc := createTestGoalCategory("update goal category", uuid.MustParse(userId))
	reqBody := map[string]any{"title": "updated title", "xp_per_goal": -1}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), reqBody)
	require.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteGoalCategoryById(t *testing.T) {
	cat := createTestGoalCategory("delete goal category", uuid.MustParse(userId))

	res, err := buildAndSendRequest("DELETE", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, cat.Id), nil)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeleteGoalCategoryByIdNotAuthorized(t *testing.T) {
	gc := createTestGoalCategory("delete goal category", uuid.MustParse(userId))
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+accessToken+"1")
	res, err := http.DefaultClient.Do(req)

	require.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestCreateGoal(t *testing.T) {
	cat := createTestGoalCategory("create goal", uuid.MustParse(userId))
	reqBody := map[string]any{
		"title": "goal title", "description": "goal description", "category_id": cat.Id,
		"user_id": userId,
	}
	res, err := buildAndSendRequest("POST", BASE_URL+"/api/goals", reqBody)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := UnmarshalServerResponse[entities.Goal](res)
	assert.Nil(t, err)
	assert.Equal(t, cat.Id.String(), resBody.Data.CategoryId.String())
	assert.Equal(t, userId, resBody.Data.UserId.String())
	assert.Equal(t, "goal title", resBody.Data.Title)
	assert.Equal(t, "goal description", resBody.Data.Description)
}

func TestCreateGoalInvalidFields(t *testing.T) {
	reqBody := map[string]any{
		"title": "goal title", "description": "goal description", "category_id": "not a uuid",
	}
	res, err := buildAndSendRequest("POST", BASE_URL+"/api/goals", reqBody)
	require.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateGoalById(t *testing.T) {
	cat := createTestGoalCategory("update goal", uuid.MustParse(userId))
	goal := createTestGoal("goal title", "goal description", cat.Id, uuid.MustParse(userId))
	reqBody := map[string]any{"title": "updated title", "description": "updated description"}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/goals/%s", BASE_URL, goal.Id), reqBody)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := UnmarshalServerResponse[entities.Goal](res)
	require.Nil(t, err)
	assert.Equal(t, "updated title", resBody.Data.Title)
	assert.Equal(t, "updated description", resBody.Data.Description)
}

func TestUserCreatedEvent(t *testing.T) {
	reqBody := handler.SignupRequest{Email: "user2@mail.com", Password: "password123!"}
	stringifiedBody, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	res, err := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()

	serverResponse, err := UnmarshalServerResponse[entities.UserDTO](res)
	assert.Nil(t, err)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/goals/categories", BASE_URL), nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+serverResponse.Data.AccessToken)
	res, err = http.DefaultClient.Do(req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := UnmarshalServerResponse[[]*entities.GoalCategory](res)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(resBody.Data))
}

func TestGoalCategoryCreatedEvent(t *testing.T) {
	// create a goal category
	body := map[string]any{"title": "testing create category event", "xp_per_goal": 100}
	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	res, err := buildAndSendRequest("POST", url, body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	responseData, err := UnmarshalServerResponse[entities.GoalCategory](res)
	gc := responseData.Data
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// check if the default goal was created
	res, err = buildAndSendRequest("GET", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), nil)
	assert.Nil(t, err)
	serverResponse, err := UnmarshalServerResponse[*entities.GoalCategory](res)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(serverResponse.Data.Goals))
}

// utility functions
func buildAndSendRequest(method, url string, body map[string]any) (*http.Response, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func printErrResponse(res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var errRes responses.APIError
	err = json.Unmarshal(body, &errRes)
	if err != nil {
		return err
	}

	return nil
}

func printSuccessResponse[T any](res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var successRes responses.ServerResponse[T]
	err = json.Unmarshal(body, &successRes)
	if err != nil {
		return err
	}
	return nil
}

func createTestGoalCategory(title string, userId uuid.UUID) *entities.GoalCategory {
	query := `INSERT INTO goal_categories (title, xp_per_goal, user_id) 
  VALUES ($1, $2, $3) RETURNING *`
	var goalCategory entities.GoalCategory
	dbx.QueryRowx(query, title, 50, userId).StructScan(&goalCategory)
	return &goalCategory
}

func createTestGoal(title, description string, categoryId, userId uuid.UUID) *entities.Goal {
	query := `INSERT INTO goals (title, description, user_id, category_id)
  VALUES ($1, $2, $3, $4) RETURNING *`
	var goal entities.Goal
	dbx.QueryRowx(query, title, description, userId, categoryId).StructScan(&goal)
	return &goal
}

func UnmarshalServerResponse[T any](res *http.Response) (responses.ServerResponse[T], error) {
	var serverResponse responses.ServerResponse[T]

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return serverResponse, err
	}

	err = json.Unmarshal(body, &serverResponse)
	return serverResponse, err
}
