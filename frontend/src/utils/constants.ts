export const API_BASE = "http://localhost:8080/api";
export const WS_BASE = "ws://localhost:8080/api/ws";

export const events = {
  USER_CREATED: "user_created",
  GOAL_CREATED: "goal_created",
  GOAL_UPDATED: "goal_updated",
  USER_UPDATED: "user_updated",
  GOAL_CATEGORY_CREATED: "goal_category_created",
  DEFAULT_GOAL_CREATED: "default_goal_created",
  SSE_CONNECTED: "sse_connected",
} as const;

export const http = {
  StatusUnauthorized: 401,
  StatusBadRequest: 400,
  StatusNotFound: 404,
  MethodPatch: "PATCH",
  MethodPost: "POST",
  MethodDelete: "DELETE",
  MethodGet: "GET",
  MethodPut: "PUT",
} as const;
