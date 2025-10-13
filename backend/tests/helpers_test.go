package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	sqlcdb "goalify/internal/db/generated"
	"goalify/internal/entities"
	"goalify/internal/users/handler"
	"goalify/internal/responses"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// HTTP helpers
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

func unmarshalResponse[T any](res *http.Response) (T, error) {
	var response T

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	return response, err
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

// User helpers
func createUser(email, password string) *entities.UserDTO {
	reqBody := handler.SignupRequest{Email: email, Password: password, ConfirmPassword: password}
	stringifiedBody, _ := json.Marshal(reqBody)
	res, _ := http.Post(BASE_URL+"/api/users/signup", "application/json", bytes.NewReader(stringifiedBody))

	servRes, _ := unmarshalResponse[*entities.UserDTO](res)
	return servRes
}

func getUserById(id string) (*entities.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	u, err := queries.GetUserById(context.Background(), pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, err
	}

	return &entities.User{
		Id:                 uuid.UUID(u.ID.Bytes),
		Email:              u.Email,
		Password:           u.Password,
		Xp:                 int(u.Xp.Int32),
		LevelId:            int(u.LevelID.Int32),
		CashAvailable:      int(u.CashAvailable.Int32),
		RefreshToken:       uuid.UUID(u.RefreshToken.Bytes),
		RefreshTokenExpiry: u.RefreshTokenExpiry.Time,
		CreatedAt:          u.CreatedAt.Time,
		UpdatedAt:          u.UpdatedAt.Time,
	}, nil
}

// Goal helpers
func createTestGoalCategory(title string, userId uuid.UUID) *entities.GoalCategory {
	params := sqlcdb.CreateGoalCategoryParams{
		Title:     title,
		XpPerGoal: 100,
		UserID:    pgtype.UUID{Bytes: userId, Valid: true},
	}
	gc, err := queries.CreateGoalCategory(context.Background(), params)
	if err != nil {
		panic(err)
	}
	return &entities.GoalCategory{
		Id:          uuid.UUID(gc.ID.Bytes),
		Title:       gc.Title,
		Xp_per_goal: int(gc.XpPerGoal),
		UserId:      uuid.UUID(gc.UserID.Bytes),
		CreatedAt:   gc.CreatedAt.Time,
		UpdatedAt:   gc.UpdatedAt.Time,
		Goals:       []*entities.Goal{},
	}
}

func createTestGoal(title, description string, categoryId, userId uuid.UUID) *entities.Goal {
	params := sqlcdb.CreateGoalParams{
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		UserID:      pgtype.UUID{Bytes: userId, Valid: true},
		CategoryID:  pgtype.UUID{Bytes: categoryId, Valid: true},
	}
	g, err := queries.CreateGoal(context.Background(), params)
	if err != nil {
		panic(err)
	}
	return &entities.Goal{
		Id:          uuid.UUID(g.ID.Bytes),
		Title:       g.Title,
		Description: g.Description.String,
		Status:      string(g.Status.GoalStatus),
		CategoryId:  uuid.UUID(g.CategoryID.Bytes),
		UserId:      uuid.UUID(g.UserID.Bytes),
		CreatedAt:   g.CreatedAt.Time,
		UpdatedAt:   g.UpdatedAt.Time,
	}
}
