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
	"goalify/testsetup"
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
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const BASE_URL = "http://localhost:8080"

var (
	dbx           *sqlx.DB
	configService *config.ConfigService
	pgContainer   *postgres.PostgresContainer
)

func setup(ctx context.Context) {
	var err error
	configService = config.NewConfigService(options.None[string]())
	configService.SetEnv(config.ENV, "test")

	pgContainer, err = testsetup.GetPgContainer()
	if err != nil {
		panic(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	dbx, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	go main.Run()
	time.Sleep(50 * time.Millisecond)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// start server in a goroutine
	setup(ctx)
	code := m.Run()

	var err error
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	os.Exit(code)
}

func TestHealth(t *testing.T) {
	t.Parallel()
	res, err := http.Get(BASE_URL + "/health")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

/* Users Domain Tests
* Testing Resource: /api/users
 */

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

	assert.Nil(t, err)
	user, err := unmarshalResponse[entities.UserDTO](res)
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

func createUser(email, password string) *entities.UserDTO {
	reqBody := handler.SignupRequest{Email: email, Password: password, ConfirmPassword: password}
	stringifiedBody, _ := json.Marshal(reqBody)
	res, _ := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))

	servRes, _ := unmarshalResponse[*entities.UserDTO](res)
	return servRes
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

/* Goal Domain Tests
* Testing Resource: /api/goals
 */

func TestGoalCategoryCreate(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	reqBody := map[string]any{"title": "goal cat", "xp_per_goal": 100}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/goals/categories", BASE_URL), reqBody, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	gc, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, "goal cat", gc.Title)
	assert.Equal(t, 100, gc.Xp_per_goal)
	assert.NotNil(t, gc.Goals)
	assert.Equal(t, userDto.Id.String(), gc.UserId.String())
}

func TestGoalCategoryCreateInvalidFields(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	reqBody := map[string]any{"title": "goal cat", "xp_per_goal": -10}
	res, err := buildAndSendRequest("POST", fmt.Sprintf("%s/api/goals/categories", BASE_URL), reqBody, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}

func TestGetGoalCategories(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	for i := 0; i < 5; i++ {
		createTestGoalCategory("create goal category", userDto.Id)
		time.Sleep(20 * time.Millisecond)
	}
	res, err := buildAndSendRequest("GET", fmt.Sprintf("%s/api/goals/categories", BASE_URL), nil, userDto.AccessToken)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := unmarshalResponse[responses.ServerResponse[[]*entities.GoalCategory]](res)
	assert.Nil(t, err)

	// expecting 5 inserted category + 1 default from user_created event
	assert.Equal(t, 6, len(resBody.Data))
}

func TestGetGoalCategoryById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	gc := createTestGoalCategory("create goal category", userDto.Id)
	url := fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)

	res, err := buildAndSendRequest("GET", url, nil, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, gc.Id, resBody.Id)
	assert.Equal(t, gc.UserId, resBody.UserId)
}

func TestUpdateGoalCategoryById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	gc := createTestGoalCategory("update goal category", userDto.Id)

	url := fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)
	reqBody := map[string]any{"title": "updated title", "xp_per_goal": 69}
	res, err := buildAndSendRequest("PUT", url, reqBody, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, "updated title", resBody.Title)
	assert.Equal(t, 69, resBody.Xp_per_goal)
}

func TestUpdateGoalCategoryByIdInvalidFields(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	gc := createTestGoalCategory("update goal category", userDto.Id)
	reqBody := map[string]any{"title": "updated title", "xp_per_goal": -1}

	url := fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)
	res, err := buildAndSendRequest("PUT", url, reqBody, userDto.AccessToken)
	require.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}

func TestDeleteGoalCategoryById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	cat := createTestGoalCategory("delete goal category", userDto.Id)

	res, err := buildAndSendRequest("DELETE", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, cat.Id), nil, userDto.AccessToken)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeleteGoalCategoryByIdNotAuthorized(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	userId := userDto.Id
	gc := createTestGoalCategory("delete goal category", userId)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id), nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+userDto.AccessToken+"1")
	res, err := http.DefaultClient.Do(req)

	require.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestCreateGoal(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	cat := createTestGoalCategory("create goal", userDto.Id)
	reqBody := map[string]any{
		"title": "goal title", "description": "goal description", "category_id": cat.Id,
		"user_id": userDto.Id,
	}
	res, err := buildAndSendRequest("POST", BASE_URL+"/api/goals", reqBody, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	resBody, err := unmarshalResponse[entities.Goal](res)
	assert.Nil(t, err)
	assert.Equal(t, cat.Id.String(), resBody.CategoryId.String())
	assert.Equal(t, userDto.Id, resBody.UserId)
	assert.Equal(t, "goal title", resBody.Title)
	assert.Equal(t, "goal description", resBody.Description)
}

func TestCreateGoalInvalidFields(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	reqBody := map[string]any{
		"title": "goal title", "description": "goal description", "category_id": "not a uuid",
	}
	res, err := buildAndSendRequest("POST", BASE_URL+"/api/goals", reqBody, userDto.AccessToken)
	require.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateGoalById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	cat := createTestGoalCategory("update goal", userDto.Id)
	goal := createTestGoal("goal title", "goal description", cat.Id, userDto.Id)
	reqBody := map[string]any{"title": "updated title", "description": "updated description"}
	res, err := buildAndSendRequest("PUT", fmt.Sprintf("%s/api/goals/%s", BASE_URL, goal.Id), reqBody, userDto.AccessToken)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := unmarshalResponse[entities.Goal](res)
	require.Nil(t, err)
	assert.Equal(t, "updated title", resBody.Title)
	assert.Equal(t, "updated description", resBody.Description)
}

func TestUserCreatedEvent(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	var resBody responses.ServerResponse[[]*entities.GoalCategory]
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/goals/categories", BASE_URL), nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+userDto.AccessToken)

	for i := 0; i < 5; i++ {
		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err = unmarshalResponse[responses.ServerResponse[[]*entities.GoalCategory]](res)
		assert.Nil(t, err)
		if len(resBody.Data) == 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	assert.Equal(t, 1, len(resBody.Data))
}

func TestGoalCategoryCreatedEvent(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	// create a goal category
	body := map[string]any{"title": "testing create category event", "xp_per_goal": 100}
	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	res, err := buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	gc, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// check if the default goal was created

	var serverResponse *entities.GoalCategory
	url = fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)
	for i := 0; i < 5; i++ {
		res, err = buildAndSendRequest("GET", url, nil, userDto.AccessToken)
		assert.Nil(t, err)
		serverResponse, err = unmarshalResponse[*entities.GoalCategory](res)
		assert.Nil(t, err)
		if len(serverResponse.Goals) == 1 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	assert.Equal(t, 1, len(serverResponse.Goals))
}

func TestDeleteGoal(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	// create category and goal to be deleted
	body := map[string]any{
		"title":       "some title",
		"xp_per_goal": 50,
	}
	url := fmt.Sprintf("%s/api/goals/categories", BASE_URL)
	res, err := buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	gc, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)

	url = fmt.Sprintf("%s/api/goals", BASE_URL)
	body = map[string]any{
		"title":       "some title",
		"description": "some description",
		"category_id": gc.Id,
	}
	assert.Nil(t, err)
	res, err = buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	goal, err := unmarshalResponse[entities.Goal](res)
	assert.Nil(t, err)

	// rerun requests until goal category created event triggers
	url = fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, gc.Id)
	for i := 0; i < 5; i++ {
		res, _ = buildAndSendRequest("GET", url, nil, userDto.AccessToken)
		gc, _ = unmarshalResponse[entities.GoalCategory](res)
		if len(gc.Goals) == 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, 2, len(gc.Goals))

	deleteUrl := fmt.Sprintf("%s/api/goals/%s", BASE_URL, goal.Id)
	res, err = buildAndSendRequest("DELETE", deleteUrl, nil, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// rerun request to get goal category by id
	res, err = buildAndSendRequest("GET", url, nil, userDto.AccessToken)
	assert.Nil(t, err)
	gc, err = unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(gc.Goals))
}

func TestDeleteGoalNotFound(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	deleteUrl := fmt.Sprintf("%s/api/goals/%s", BASE_URL, uuid.New())
	res, err := buildAndSendRequest("DELETE", deleteUrl, nil, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUserLevelUpEvent(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	cat := createTestGoalCategory("create goal category", userDto.Id)
	goal := createTestGoal("create goal", "desc", cat.Id, userDto.Id)

	reqBody := map[string]any{"status": "complete"}
	url := fmt.Sprintf("%s/api/goals/%s", BASE_URL, goal.Id)
	res, err := buildAndSendRequest("PUT", url, reqBody, userDto.AccessToken)
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// rerun requests until goal category created event triggers
	url = fmt.Sprintf("%s/api/goals/categories/%s", BASE_URL, cat.Id)
	var user *entities.User
	for i := 0; i < 5; i++ {
		user, err = getUserById(userDto.Id.String())
		if err != nil {
			t.Log(err)
			break
		}
		if user.LevelId == 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, 2, user.LevelId)
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

// utility functions
func buildAndSendRequest(method, url string, body map[string]any, accessToken string) (*http.Response, error) {
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

func createTestGoalCategory(title string, userId uuid.UUID) *entities.GoalCategory {
	query := `INSERT INTO goal_categories (title, xp_per_goal, user_id) 
  VALUES ($1, $2, $3) RETURNING *`
	var goalCategory entities.GoalCategory
	dbx.QueryRowx(query, title, 100, userId).StructScan(&goalCategory)
	return &goalCategory
}

func createTestGoal(title, description string, categoryId, userId uuid.UUID) *entities.Goal {
	query := `INSERT INTO goals (title, description, user_id, category_id)
  VALUES ($1, $2, $3, $4) RETURNING *`
	var goal entities.Goal
	dbx.QueryRowx(query, title, description, userId, categoryId).StructScan(&goal)
	return &goal
}

func getUserById(id string) (*entities.User, error) {
	var user entities.User
	err := dbx.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func unmarshalResponse[T any](res *http.Response) (T, error) {
	var response T

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
}
