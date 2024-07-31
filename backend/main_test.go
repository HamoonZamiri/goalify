package main_test

import (
	"bytes"
	"context"
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
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const MIGRATION_STR string = `
CREATE TABLE levels  (
    id SERIAL PRIMARY KEY,
    level_up_xp INTEGER NOT NULL,
    cash_reward INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    xp INTEGER DEFAULT 0,
    level_id SERIAL REFERENCES levels(id),
    cash_available INTEGER DEFAULT 0,
    refresh_token UUID DEFAULT gen_random_uuid(),
    refresh_token_expiry TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE goal_categories  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    xp_per_goal INTEGER NOT NULL,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TYPE goal_status AS ENUM ('complete', 'not_complete');

CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT '',
    user_id UUID REFERENCES users(id),
    category_id UUID REFERENCES goal_categories(id) ON DELETE CASCADE,
    status goal_status DEFAULT 'not_complete',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Insert Default Levels
INSERT INTO levels (id, level_up_xp, cash_reward) VALUES (1, 100, 10);
`

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
	ctx := context.Background()

	var err error
	configService = config.NewConfigService(options.None[string]())
	configService.SetEnv("ENV", "test")
	dbName := configService.MustGetEnv("TEST_DB_NAME")
	dbUser := configService.MustGetEnv("DB_USER")
	dbPassword := configService.MustGetEnv("DB_PASSWORD")

	pgContainer, err := postgres.Run(ctx, "docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}
	configService.SetEnv("TEST_DB_CONN_STRING", connStr)

	dbx, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	dbx.MustExec(MIGRATION_STR)

	// start server in a goroutine
	setup()
	code := m.Run()

	cleanup := `DELETE from goals; DELETE FROM goal_categories; DELETE from users;`
	dbx.MustExec(cleanup)
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
	assert.Nil(t, err)
	serverResponse, err := UnmarshalServerResponse[entities.UserDTO](res)
	assert.Nil(t, err)
	refreshToken = serverResponse.Data.RefreshToken.String()
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

	var resBody responses.ServerResponse[[]*entities.GoalCategory]
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/goals/categories", BASE_URL), nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+serverResponse.Data.AccessToken)
	for i := 0; i < 5; i++ {
		res, err = http.DefaultClient.Do(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err = UnmarshalServerResponse[[]*entities.GoalCategory](res)
		assert.Nil(t, err)
		if len(resBody.Data) == 1 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

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

	var serverResponse responses.ServerResponse[*entities.GoalCategory]
	for i := 0; i < 5; i++ {
		res, err = buildAndSendRequest("GET", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), nil)
		assert.Nil(t, err)
		serverResponse, err = UnmarshalServerResponse[*entities.GoalCategory](res)
		assert.Nil(t, err)
		if len(serverResponse.Data.Goals) == 1 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	assert.Equal(t, 1, len(serverResponse.Data.Goals))
}

func TestDeleteGoal(t *testing.T) {
	// create category and goal to be deleted
	body := map[string]any{
		"title":       "some title",
		"xp_per_goal": 50,
	}
	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	res, err := buildAndSendRequest("POST", url, body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	gc, err := UnmarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)

	url = fmt.Sprintf("%s/api/goals", BASE_URL)
	body = map[string]any{
		"title":       "some title",
		"description": "some description",
		"category_id": gc.Data.Id,
	}
	res, err = buildAndSendRequest("POST", url, body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	goal, err := UnmarshalServerResponse[entities.Goal](res)
	assert.Nil(t, err)

	// rerun requests until goal category created event triggers
	url = fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Data.Id)
	for i := 0; i < 5; i++ {
		res, _ = buildAndSendRequest("GET", url, nil)
		gc, _ = UnmarshalServerResponse[entities.GoalCategory](res)
		if len(gc.Data.Goals) == 2 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	assert.Equal(t, 2, len(gc.Data.Goals))

	deleteUrl := fmt.Sprintf("%s/api/goals/%s", BASE_URL, goal.Data.Id)
	res, err = buildAndSendRequest("DELETE", deleteUrl, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// rerun request to get goal category by id
	res, err = buildAndSendRequest("GET", url, nil)
	assert.Nil(t, err)
	gc, err = UnmarshalServerResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(gc.Data.Goals))
}

func TestDeleteGoalNotFound(t *testing.T) {
	deleteUrl := fmt.Sprintf("%s/api/goals/%s", BASE_URL, uuid.New())
	res, err := buildAndSendRequest("DELETE", deleteUrl, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
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
