package user

import (
	"context"
	"net/http"
	"os"
	"proj/helper"

	"github.com/golang-jwt/jwt/v5"
)

// TODO add a refresh token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieToken, err := r.Cookie("token")
		if err != nil {
			helper.SendError(w, http.StatusUnauthorized, "unauthorized")

			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(cookieToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil //TODO keep secret in config
		})
		if err != nil || !token.Valid {
			helper.SendError(w, http.StatusBadRequest, "token is not valid")

			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
