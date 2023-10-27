package userModel

import "github.com/golang-jwt/jwt/v5"

type UserBody struct {
	Username       string `json:"username"`
	Fullname       string `json:"fullname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Number     string `json:"number"`
	Avatar     string `json:"avatar"`
	Location    string `json:"location"`
	Coordinates    string `json:"coordinates"`
	Created_on string `json:"created_on"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
