package middlewares

import (
	"context"
	"net/http"
	"server/internal/app/adapters/http-adapter/codec"
	"server/internal/app/adapters/http-adapter/constants"
	errorMapper "server/internal/app/adapters/http-adapter/errors/user-usecase"
)

type authService interface {
	Authenticate(token string) (int64, error)
}

func JWTMiddleware(service authService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt := r.Header.Get("Authorization")
			if jwt == "" {
				http.Error(w, "JWT token not found", http.StatusUnauthorized)
				return
			}
			userID, err := service.Authenticate(jwt)
			if err != nil {
				status, message := errorMapper.Process(err)
				codec.WriteJSON(w, status, message)
				return
			}
			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
