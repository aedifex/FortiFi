package middleware

import (
	"context"
	"net/http"

	"github.com/aedifex/FortiFi/pkg/utils"
	"go.uber.org/zap"
)

type contextKey string
const UserIdContextKey contextKey = "userId"

// Middleware to check jwt for protected endpoints
func Auth(key string, logger *zap.SugaredLogger, next http.HandlerFunc) http.HandlerFunc {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		id, err := utils.GetJwtId(key, request.Header.Get("Authorization"))
		if err != nil {
			logger.Infof("failed to get jwt id: %s", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), UserIdContextKey , id)
		next.ServeHTTP(writer,request.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}