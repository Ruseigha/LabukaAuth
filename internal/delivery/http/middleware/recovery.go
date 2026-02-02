package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery recovers from panics and returns 500 error
// WHY: Prevent server crashes, log stack traces
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic with stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Return 500 to client
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error","message":"an unexpected error occurred","code":500}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
