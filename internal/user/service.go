package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"proj/helper"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	UserRepo UserRepository
}

func (service UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	user, err := service.UserRepo.GetUserById(userId)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	helper.SendJSON(w, http.StatusOK, user)
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

	userId, err := service.UserRepo.CheckUserCredentials(credentials)
	switch {
	case err == sql.ErrNoRows:
		helper.SendError(w, http.StatusBadRequest, "user is not exists")

		return
	case err != nil:
		helper.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	expirationTime := time.Now().Add(12 * time.Hour)

	claims := &Claims{
		UserName: credentials.Name,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) //TODO keep secret in config

	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	helper.SendJSON(w, http.StatusOK, nil)
}

func (service UserService) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	helper.SendJSON(w, http.StatusOK, nil)
}
