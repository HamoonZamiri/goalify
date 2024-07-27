export const API_BASE = "http://localhost:8080/api";

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
