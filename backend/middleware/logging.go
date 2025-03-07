package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusCodeWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (scw *statusCodeWriter) WriteHeader(statusCode int) {
	scw.statusCode = statusCode
	scw.w.WriteHeader(statusCode)
}

func (scw *statusCodeWriter) Header() http.Header {
	return scw.w.Header()
}

func (scw *statusCodeWriter) Write(b []byte) (int, error) {
	return scw.w.Write(b)
}

func (scw *statusCodeWriter) Flush() {
	scw.w.(http.Flusher).Flush()
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusWriter := &statusCodeWriter{w: w, statusCode: http.StatusOK}
		next.ServeHTTP(statusWriter, r)
		slog.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", statusWriter.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
