package userModel

import "github.com/golang-jwt/jwt/v5"

type UserBody struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Number     string `json:"number"`
	Avatar     string `json:"avatar"`
	Address    string `json:"address"`
	Created_on string `json:"created_on"`
}

type Claims struct {
	Username string `json:"username"`

	jwt.RegisteredClaims
}
