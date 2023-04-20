package user

import (
	"encoding/json"
	"net/http"
	"proj/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	UserRepo UserRepository
}

func (service UserService) Ping(w http.ResponseWriter, r *http.Request) {

}

func (service UserService) Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	validate := validator.New()

	err = validate.Struct(credentials)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, "bad request")

		return
	}

	isExists, err := service.UserRepo.CheckIfUserExists(credentials)
	switch {
	case err != nil:
		helper.SendError(w, http.StatusBadRequest, err.Error())

		return
	case !isExists:
		helper.SendError(w, http.StatusBadRequest, "user is not exists")

		return
	}

	claims := &Claims{
		Username: credentials.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret")) //TODO keep secret in config

	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	helper.SendJSON(w, http.StatusOK, tokenString)
}
