package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// Middleware to log all requests to the server
func Logging(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			logger.Infof("Incoming request from: %s", request.RemoteAddr)
			next.ServeHTTP(writer, request) // Call the next handler
		})
	}
}