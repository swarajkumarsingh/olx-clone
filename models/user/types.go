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
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
}

type User struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Fullname    string `json:"fullname" db:"fullname"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Phone       string `json:"phone" db:"phone"`
	Avatar      string `json:"avatar" db:"avatar"`
	Location    string `json:"location" db:"location"`
	Coordinates string `json:"coordinates" db:"coordinates"`
	Created_at  string `json:"created_on" db:"created_at"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
