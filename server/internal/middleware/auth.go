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
		tokenHeader := request.Header.Get("Authorization")
		signedToken, err := utils.ExtractBearer(tokenHeader)
		if err != nil {
			logger.Infof("failed to extract jwt id: %s", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := utils.GetJwtSubject(key, signedToken)
		if err != nil {
			logger.Infof("failed to get id from jwt: %s", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), UserIdContextKey , id)
		next.ServeHTTP(writer,request.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}