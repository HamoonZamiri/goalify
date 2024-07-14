package middleware

import (
	"log"
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusWriter := &statusCodeWriter{w: w, statusCode: http.StatusOK}
		next.ServeHTTP(statusWriter, r)
		log.Println(statusWriter.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
