package tests

import (
	"fmt"
	"goalify/internal/entities"
	"goalify/internal/responses"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* Goal Category Domain Tests
* Testing Resource: /api/goals/categories
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
