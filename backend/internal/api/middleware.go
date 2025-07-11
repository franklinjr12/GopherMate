package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// LoggingMiddleware logs the route and payload of every request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			// Restore the io.ReadCloser to its original state
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		var logMessage string
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			logMessage = "Payload: " + string(bodyBytes)
		}
		log.Printf("%s %s %s", r.Method, r.URL.Path, logMessage)
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Allow requests from your frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
