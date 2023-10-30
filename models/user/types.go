package userModel

import "github.com/golang-jwt/jwt/v5"

type UserBody struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
	Created_at  string `json:"created_at"`
}

type UserUpdateBody struct {
	Username        string `json:"username"`
	Fullname        string `json:"fullname"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar"`
	Location        string `json:"location"`
	Coordinates     string `json:"coordinates"`
}

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
	Created_at  string `json:"created_on"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
