package logger

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// wraps our RW in a struct to manage logging
func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// Status returns our rw's status code
func (rw *responseWriter) Status() int {
	return rw.status
}

// WriteHeader turns the wroteHeader property true
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware wraps our router to get a handle on
func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Case where there's a panic
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		fmt.Println(
			"status", wrapped.status,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	}

	return http.HandlerFunc(fn)
}
