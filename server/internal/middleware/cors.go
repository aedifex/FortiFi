package middleware

import (
	"net/http"
)

func CORSMiddleware(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", origin)
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")	
			
			if request.Method == "OPTIONS" {
				writer.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(writer, request) // Call the next handler
		})
	}
}