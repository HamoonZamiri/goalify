// Package responses contains utilities for sending REST responses to clients
package responses

import (
	"goalify/pkg/jsonutil"
	"log/slog"
	"net/http"
)

type (
	Object                string
	Type                  string
	ServerResponse[T any] struct {
		Data     T       `json:"data"`
		Type     *Type   `json:"type,omitempty"`
		HasMore  *bool   `json:"has_more,omitempty"`
		NextPage *string `json:"next_page,omitempty"`
		Object   Object  `json:"object"`
	}
)

const (
	ObjectGoalCategory Object = "goal_category"
	ObjectUser         Object = "user"
	ObjectGoal         Object = "goal"
	ObjectList         Object = "list"
)

func SendResponse[T any | map[string]any](
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data T,
) {
	if err := jsonutil.Encode(w, r, status, data); err != nil {
		slog.Error("responses.SendResponse: jsonutil.Encode: ", "err", err)
		SendAPIError(w, r, http.StatusInternalServerError, ErrInternalServer.Error(), nil)
	}
}
