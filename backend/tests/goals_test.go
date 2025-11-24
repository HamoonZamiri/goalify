package tests

import (
	"fmt"
	"goalify/internal/entities"
	"goalify/internal/responses"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* Goal Domain Tests
* Testing Resource: /api/goals
 */

func TestCreateGoal(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	cat := createTestGoalCategory("create goal", userDto.ID)
	reqBody := map[string]any{
		"title": "goal title", "description": "goal description", "category_id": cat.ID,
		"user_id": userDto.ID,
	}
	res, err := buildAndSendRequest("POST", BaseURL+"/api/goals", reqBody, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	resBody, err := unmarshalResponse[entities.Goal](res)
	assert.Nil(t, err)
	assert.Equal(t, cat.ID.String(), resBody.CategoryID.String())
	assert.Equal(t, userDto.ID, resBody.UserID)
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
	res, err := buildAndSendRequest("POST", BaseURL+"/api/goals", reqBody, userDto.AccessToken)
	require.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateGoalById(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	cat := createTestGoalCategory("update goal", userDto.ID)
	goal := createTestGoal("goal title", "goal description", cat.ID, userDto.ID)
	reqBody := map[string]any{"title": "updated title", "description": "updated description"}
	res, err := buildAndSendRequest(
		"PUT",
		fmt.Sprintf("%s/api/goals/%s", BaseURL, goal.ID),
		reqBody,
		userDto.AccessToken,
	)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := unmarshalResponse[entities.Goal](res)
	require.Nil(t, err)
	assert.Equal(t, "updated title", resBody.Title)
	assert.Equal(t, "updated description", resBody.Description)
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
	url := fmt.Sprintf("%s/api/goals/categories", BaseURL)
	res, err := buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	gc, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)

	url = fmt.Sprintf("%s/api/goals", BaseURL)
	body = map[string]any{
		"title":       "some title",
		"description": "some description",
		"category_id": gc.ID,
	}
	assert.Nil(t, err)
	res, err = buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	goal, err := unmarshalResponse[entities.Goal](res)
	assert.Nil(t, err)

	// rerun requests until goal category created event triggers
	url = fmt.Sprintf("%s/api/goals/categories/%s", BaseURL, gc.ID)
	for range 5 {
		res, _ = buildAndSendRequest("GET", url, nil, userDto.AccessToken)
		gc, _ = unmarshalResponse[entities.GoalCategory](res)
		if len(gc.Goals) == 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, 2, len(gc.Goals))

	deleteURL := fmt.Sprintf("%s/api/goals/%s", BaseURL, goal.ID)
	res, err = buildAndSendRequest("DELETE", deleteURL, nil, userDto.AccessToken)
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

	deleteURL := fmt.Sprintf("%s/api/goals/%s", BaseURL, uuid.New())
	res, err := buildAndSendRequest("DELETE", deleteURL, nil, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

/* Event Tests
* Testing event-driven behavior
 */

func TestUserLevelUpEvent(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")
	cat := createTestGoalCategory("create goal category", userDto.ID)
	goal := createTestGoal("create goal", "desc", cat.ID, userDto.ID)

	reqBody := map[string]any{"status": "complete"}
	url := fmt.Sprintf("%s/api/goals/%s", BaseURL, goal.ID)
	res, err := buildAndSendRequest("PUT", url, reqBody, userDto.AccessToken)
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var user *entities.User
	for i := 0; i < 5; i++ {
		user, err = getUserByID(userDto.ID.String())
		if err != nil {
			t.Log(err)
			break
		}
		if user.LevelID == 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, 2, user.LevelID)
}

func TestUserCreatedEvent(t *testing.T) {
	t.Parallel()

	email := t.Name() + "@mail.com"
	userDto := createUser(email, "password123!")

	var resBody responses.ServerResponse[[]*entities.GoalCategory]
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/goals/categories", BaseURL), nil)
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
	url := fmt.Sprintf("%s/api/goals/categories", BaseURL)
	res, err := buildAndSendRequest("POST", url, body, userDto.AccessToken)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	gc, err := unmarshalResponse[entities.GoalCategory](res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// check if the default goal was created

	var serverResponse *entities.GoalCategory
	url = fmt.Sprintf("%s/api/goals/categories/%s", BaseURL, gc.ID)
	for range 5 {
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
