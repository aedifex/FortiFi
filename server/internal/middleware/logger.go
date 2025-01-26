package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// Middleware to log all requests to the server
func Logging(logger *zap.SugaredLogger, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, request *http.Request) {
		logger.Infof("Incoming request from: %s", request.RemoteAddr)
		next.ServeHTTP(w,request)
	}
	return http.HandlerFunc(fn)
}