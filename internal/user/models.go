package user

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id   int    `json:"id" validate:"required,numeric"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,numeric"`
}

type Credentials struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
