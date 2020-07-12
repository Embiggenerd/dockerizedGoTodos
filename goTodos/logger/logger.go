package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

type requestLog struct {
	Status   int
	Method   string
	Path     string
	Duration time.Duration
}

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

// WriteHeader turns the wroteHeader property true to prevent extra write
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware wraps our router to get a handle on requests. If a logs file
// exists, it adds request log to json array. Else, it creates new file and logs request
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

		request := requestLog{
			Status:   wrapped.status,
			Method:   r.Method,
			Path:     r.URL.EscapedPath(),
			Duration: time.Since(start),
		}

		logsFilepath := filepath.Join("logger", "logs.json")

		var logsArr []interface{}

		if fileExists(logsFilepath) {

			logs, _ := ioutil.ReadFile(logsFilepath)

			json.Unmarshal(logs, &logsArr)

			logsArr = append(logsArr, request)

			logsArrJSON, _ := json.MarshalIndent(&logsArr, "", " ")

			ioutil.WriteFile(logsFilepath, logsArrJSON, 0644)

		} else {

			logsArr = append(logsArr, request)

			logsArrJSON, _ := json.MarshalIndent(&logsArr, "", " ")

			ioutil.WriteFile(logsFilepath, logsArrJSON, 0644)
		}
	}

	return http.HandlerFunc(fn)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
