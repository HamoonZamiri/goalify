package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type statusCodeWriter struct {
	w          http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

var redactedKeys = map[string]bool{
	"password":             true,
	"confirm_password":     true,
	"access_token":         true,
	"refresh_token":        true,
	"refresh_token_expiry": true,
}

func (scw *statusCodeWriter) WriteHeader(statusCode int) {
	scw.statusCode = statusCode
	scw.w.WriteHeader(statusCode)
}

func (scw *statusCodeWriter) Header() http.Header {
	return scw.w.Header()
}

func (scw *statusCodeWriter) Write(b []byte) (int, error) {
	scw.body.Write(b)
	return scw.w.Write(b)
}

func (scw *statusCodeWriter) Flush() {
	flusher, ok := scw.w.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func prettifyAndRedactJSON(data []byte, wasTruncated bool) any {
	var jsonData map[string]any
	defaultResponse := string(data)

	// If truncated, add indicator to plain string response
	if wasTruncated {
		defaultResponse += "... [TRUNCATED]"
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return defaultResponse
	}

	for key := range jsonData {
		if redactedKeys[key] {
			jsonData[key] = "[REDACTED]"
		}
	}

	// For valid JSON that was truncated, add indicator to map
	if wasTruncated {
		jsonData["_truncated"] = true
	}

	return jsonData
}

const maxBodySizeForLogging = 1024 // 1KB limit for logged request/response bodies

func isLongLivedConnection(path string) bool {
	return strings.Contains(path, "/events") || strings.Contains(path, "/ws")
}

func truncateBody(data []byte, maxSize int) []byte {
	if len(data) <= maxSize {
		return data
	}
	return data[:maxSize]
}

func redactURLForSSE(path, rawQuery string) string {
	// Only redact for SSE endpoint
	if !strings.Contains(path, "/events") {
		if rawQuery != "" {
			return path + "?" + rawQuery
		}
		return path
	}

	// For SSE, remove access_token from query params
	if rawQuery == "" {
		return path
	}

	params := strings.Split(rawQuery, "&")
	filtered := make([]string, 0, len(params))

	for _, param := range params {
		if !strings.HasPrefix(param, "token=") {
			filtered = append(filtered, param)
		}
	}

	if len(filtered) > 0 {
		return path + "?" + strings.Join(filtered, "&")
	}
	return path
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log long-lived connections but don't log full request/response
		if isLongLivedConnection(r.URL.Path) {
			redactedURL := redactURLForSSE(r.URL.Path, r.URL.RawQuery)
			slog.Info("long-lived connection established",
				slog.String("method", r.Method),
				slog.String("path", redactedURL),
				slog.String("type", "persistent"),
			)
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Read and log request body for POST/PUT/PATCH methods
		var prettyRequestBody any
		if r.Method == http.MethodPost || r.Method == http.MethodPut ||
			r.Method == http.MethodPatch {
			// Read the request body
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Logging: io.ReadAll: ", "err", err)
				return
			}
			r.Body.Close()

			// Create a new reader with the same content for the handler
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// Process request body for logging (truncate and redact)
			requestBodySize := len(requestBody)
			wasRequestTruncated := requestBodySize > maxBodySizeForLogging
			truncatedRequestBody := truncateBody(requestBody, maxBodySizeForLogging)
			prettyRequestBody = prettifyAndRedactJSON(truncatedRequestBody, wasRequestTruncated)
		}

		statusWriter := &statusCodeWriter{w: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		next.ServeHTTP(statusWriter, r)

		// Truncate response body before logging
		bodySize := statusWriter.body.Len()
		wasTruncated := bodySize > maxBodySizeForLogging
		truncatedBody := truncateBody(statusWriter.body.Bytes(), maxBodySizeForLogging)
		prettyResponseBody := prettifyAndRedactJSON(truncatedBody, wasTruncated)
		durationLog := fmt.Sprintf("%dms", time.Since(start).Milliseconds())

		logArgs := []any{
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
			slog.Group("response",
				slog.Int("status", statusWriter.statusCode),
				slog.String("duration", durationLog),
				slog.Any("body", prettyResponseBody),
				slog.Int("body_size", bodySize),
				slog.Bool("truncated", wasTruncated),
			),
		}

		// Add request body to logs if it was present
		if prettyRequestBody != nil {
			logArgs = append(logArgs, slog.Any("request_body", prettyRequestBody))
		}

		slog.Info("incoming request", logArgs...)
	})
}
