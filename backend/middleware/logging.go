package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type statusCodeWriter struct {
	w          http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

var redactedKeys = map[string]bool{
	"password": true,
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
	scw.w.(http.Flusher).Flush()
}

func prettifyAndRedactJSON(data []byte) interface{} {
	var jsonData map[string]interface{}
	defaultResponse := string(data)
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return defaultResponse
	}

	for key := range jsonData {
		if redactedKeys[key] {
			jsonData[key] = "[REDACTED]"
		}
	}
	return jsonData
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusWriter := &statusCodeWriter{w: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		next.ServeHTTP(statusWriter, r)
		prettyResponseBody := prettifyAndRedactJSON(statusWriter.body.Bytes())
		durationLog := fmt.Sprintf("%dms", time.Since(start).Milliseconds())

		slog.Info("incoming request",
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
			slog.Group("response",
				slog.Int("status", statusWriter.statusCode),
				slog.String("duration", durationLog),
				slog.Any("body", prettyResponseBody),
				slog.Int("body_size", statusWriter.body.Len()),
			),
		)
	})
}
