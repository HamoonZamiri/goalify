package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogging_LogsSSEConnectionEstablishment(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("event: connected\ndata: hello\n\n"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/events?token=secret123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify log was created for connection establishment
	assert.NotEmpty(t, logBuf.String(), "SSE route should log connection establishment")

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify it's a long-lived connection log
	assert.Equal(t, "long-lived connection established", logEntry["msg"])
	assert.Equal(t, "GET", logEntry["method"])
	assert.Equal(t, "persistent", logEntry["type"])

	// Verify access_token is redacted from URL
	path := logEntry["path"].(string)
	assert.NotContains(t, path, "secret123", "Access token should be redacted")
	assert.NotContains(t, path, "token=", "token param should be removed")
	assert.Contains(t, path, "/api/events", "Path should still be present")

	// Verify no request/response body in logs
	assert.NotContains(t, logEntry, "response")
	assert.NotContains(t, logEntry, "request_body")
}

func TestLogging_LogsWebSocketConnectionEstablishment(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/ws", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify log was created for connection establishment
	assert.NotEmpty(t, logBuf.String(), "WebSocket route should log connection establishment")

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "long-lived connection established", logEntry["msg"])
	assert.Equal(t, "persistent", logEntry["type"])
}

func TestLogging_LogsRegularRoutes(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	responseBody := map[string]string{"message": "success"}
	responseJSON, _ := json.Marshal(responseBody)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify log was created
	assert.NotEmpty(t, logBuf.String(), "Regular route should generate logs")

	// Parse log output
	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify log contains expected fields
	assert.Equal(t, "incoming request", logEntry["msg"])
	assert.Contains(t, logEntry, "request")
	assert.Contains(t, logEntry, "response")

	// Verify request details
	request := logEntry["request"].(map[string]any)
	assert.Equal(t, "GET", request["method"])
	assert.Equal(t, "/api/users", request["path"])

	// Verify response details
	response := logEntry["response"].(map[string]any)
	assert.Equal(t, float64(http.StatusOK), response["status"])
	assert.Contains(t, response, "duration")
}

func TestLogging_TruncatesLargeResponses(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	// Create response larger than 1KB
	largeBody := strings.Repeat("x", 2000)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(largeBody))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify client receives full response (not truncated)
	assert.Equal(t, 2000, w.Body.Len(), "Client should receive full response")
	assert.Equal(t, largeBody, w.Body.String(), "Client response should not be truncated")

	// Parse log output
	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	response := logEntry["response"].(map[string]any)

	// Verify body_size reflects actual size
	assert.Equal(t, float64(2000), response["body_size"])

	// Verify truncated flag is set
	assert.Equal(t, true, response["truncated"])

	// Verify logged body is smaller than actual
	loggedBody := response["body"].(string)
	assert.LessOrEqual(t, len(loggedBody), maxBodySizeForLogging+len("... [TRUNCATED]"),
		"Logged body should be truncated to max size")
	assert.Contains(t, loggedBody, "[TRUNCATED]", "Truncated log should have indicator")
}

func TestLogging_TruncatesLargeJSONResponses(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	// Create large JSON response (larger than 1KB)
	items := make([]map[string]any, 0)
	for i := range 50 {
		items = append(items, map[string]any{
			"id":   i,
			"name": strings.Repeat("x", 20),
			"data": strings.Repeat("y", 20),
		})
	}
	largeJSON, _ := json.Marshal(map[string]any{"items": items})

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(largeJSON)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify client receives full JSON response
	assert.Greater(t, w.Body.Len(), maxBodySizeForLogging)
	var clientJSON map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &clientJSON)
	require.NoError(t, err, "Client should receive valid JSON")

	// Parse log output
	var logEntry map[string]any
	err = json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	response := logEntry["response"].(map[string]any)

	// Verify truncated flag is set
	assert.Equal(t, true, response["truncated"])

	// If the truncated JSON is still valid, it should have _truncated marker
	if bodyMap, ok := response["body"].(map[string]any); ok {
		assert.Equal(t, true, bodyMap["_truncated"],
			"Valid JSON logs should have _truncated field")
	}
}

func TestLogging_DoesNotTruncateSmallResponses(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	smallBody := `{"message": "small response"}`

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(smallBody))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	response := logEntry["response"].(map[string]any)

	// Verify truncated flag is false
	assert.Equal(t, false, response["truncated"])

	// Verify body_size matches actual size
	assert.Equal(t, float64(len(smallBody)), response["body_size"])
}

func TestLogging_LogsRequestBodyForPOST(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	requestBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
	}
	requestJSON, _ := json.Marshal(requestBody)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body to simulate handler behavior
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, requestJSON, body, "Handler should receive full request body")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123"}`))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(requestJSON))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify request body is logged
	assert.Contains(t, logEntry, "request_body")
	requestBodyLog := logEntry["request_body"].(map[string]any)
	assert.Equal(t, "testuser", requestBodyLog["username"])
	assert.Equal(t, "test@example.com", requestBodyLog["email"])
}

func TestLogging_LogsRequestBodyForPUT(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	requestBody := map[string]string{"name": "updated"}
	requestJSON, _ := json.Marshal(requestBody)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))

	req := httptest.NewRequest(http.MethodPut, "/api/goals/123", bytes.NewBuffer(requestJSON))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify request body is logged for PUT
	assert.Contains(t, logEntry, "request_body")
}

func TestLogging_DoesNotLogRequestBodyForGET(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "test"}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var logEntry map[string]any
	err := json.Unmarshal(logBuf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify request body is NOT logged for GET
	assert.NotContains(t, logEntry, "request_body")
}

func TestLogging_RedactsSensitiveFieldsInResponseBody(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	sensitiveBody := map[string]string{
		"username":         "testuser",
		"password":         "secret123",
		"access_token":     "token123",
		"refresh_token":    "refresh123",
		"confirm_password": "confirmsecret123",
	}
	responseJSON, _ := json.Marshal(sensitiveBody)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Verify sensitive fields are redacted in response
	assert.NotContains(t, logOutput, "secret123", "Password should be redacted")
	assert.NotContains(t, logOutput, "token123", "Access token should be redacted")
	assert.NotContains(t, logOutput, "refresh123", "Refresh token should be redacted")
	assert.NotContains(t, logOutput, "confirmsecret123", "Refresh token should be redacted")
	assert.Contains(t, logOutput, "[REDACTED]", "Redacted placeholder should be present")

	// Verify non-sensitive fields are still logged
	assert.Contains(t, logOutput, "testuser", "Non-sensitive fields should be logged")
}

func TestLogging_RedactsSensitiveFieldsInRequestBody(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	requestBody := map[string]string{
		"username":         "newuser",
		"password":         "supersecret456",
		"email":            "user@example.com",
		"confirm_password": "confirmsecret123",
	}
	requestJSON, _ := json.Marshal(requestBody)

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123"}`))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/users/signup", bytes.NewBuffer(requestJSON))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	logOutput := logBuf.String()

	// Verify password is redacted in request body
	assert.NotContains(t, logOutput, "supersecret456", "Password in request should be redacted")
	assert.NotContains(
		t,
		logOutput,
		"confirmsecret123",
		"Confirm Password in request should be redacted",
	)
	assert.Contains(t, logOutput, "[REDACTED]", "Redacted placeholder should be present")

	// Verify non-sensitive fields are still logged
	assert.Contains(t, logOutput, "newuser", "Username should be logged")
	assert.Contains(t, logOutput, "user@example.com", "Email should be logged")
}

func TestIsLongLivedConnection(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/events", true},
		{"/api/ws", true},
		{"/api/users", false},
		{"/api/goals", false},
		{"/events/special", true}, // Contains /events
		{"/api/websocket", false}, // Doesn't contain /ws exactly
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isLongLivedConnection(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncateBody(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		maxSize     int
		expectedLen int
	}{
		{
			name:        "body smaller than limit",
			input:       []byte("small"),
			maxSize:     100,
			expectedLen: 5,
		},
		{
			name:        "body equal to limit",
			input:       []byte(strings.Repeat("x", 100)),
			maxSize:     100,
			expectedLen: 100,
		},
		{
			name:        "body larger than limit",
			input:       []byte(strings.Repeat("x", 200)),
			maxSize:     100,
			expectedLen: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateBody(tt.input, tt.maxSize)
			assert.Equal(t, tt.expectedLen, len(result))
			if len(tt.input) > tt.maxSize {
				assert.Equal(t, tt.input[:tt.maxSize], result)
			}
		})
	}
}

func TestRedactURLForSSE(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		rawQuery string
		expected string
	}{
		{
			name:     "SSE with access_token",
			path:     "/api/events",
			rawQuery: "token=secret123",
			expected: "/api/events",
		},
		{
			name:     "SSE with access_token and other params",
			path:     "/api/events",
			rawQuery: "token=secret123&user_id=456",
			expected: "/api/events?user_id=456",
		},
		{
			name:     "SSE with multiple params including access_token",
			path:     "/api/events",
			rawQuery: "foo=bar&token=secret123&baz=qux",
			expected: "/api/events?foo=bar&baz=qux",
		},
		{
			name:     "SSE without query params",
			path:     "/api/events",
			rawQuery: "",
			expected: "/api/events",
		},
		{
			name:     "SSE without token",
			path:     "/api/events",
			rawQuery: "user_id=123&channel=main",
			expected: "/api/events?user_id=123&channel=main",
		},
		{
			name:     "Non-SSE route with access_token (should not redact)",
			path:     "/api/users",
			rawQuery: "token=secret123",
			expected: "/api/users?token=secret123",
		},
		{
			name:     "WebSocket route (not SSE, no redaction)",
			path:     "/api/ws",
			rawQuery: "token=secret123",
			expected: "/api/ws?token=secret123",
		},
		{
			name:     "Non-SSE route without query params",
			path:     "/api/goals",
			rawQuery: "",
			expected: "/api/goals",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactURLForSSE(tt.path, tt.rawQuery)
			assert.Equal(t, tt.expected, result)
		})
	}
}
