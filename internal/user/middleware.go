package user

import (
	"fmt"
	"net/http"
	"proj/helper"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// TODO add a refresh token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer")
		if len(splitToken) != 2 {
			helper.SendError(w, http.StatusBadRequest, "token is not valid")

			return
		}

		tokenString := strings.TrimSpace(splitToken[1])
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil //TODO keep secret in config
		})
		if err != nil || !token.Valid {
			fmt.Println(err.Error())
			helper.SendError(w, http.StatusUnauthorized, "unauthorized")

			return
		}

		helper.SendJSON(w, http.StatusOK, token)
		next.ServeHTTP(w, r)
	})
}
