package routes

import (
	"expvar"
	"net/http"
	"runtime"
)

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
}

// Metrics updates program counters.
func Metrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Increment the request counter.
		m.req.Add(1)

		// Update the count for the number of active goroutines every 100 requests.
		if m.req.Value()%100 == 0 {
			m.gr.Set(int64(runtime.NumGoroutine()))
		}
	}

	return http.HandlerFunc(fn)
}
