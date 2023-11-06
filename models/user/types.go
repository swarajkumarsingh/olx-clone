package userModel

import "github.com/golang-jwt/jwt/v5"

type UserBody struct {
	Username    string `validate:"required" json:"username"`
	Fullname    string `validate:"required" json:"fullname"`
	Email       string `validate:"required" json:"email"`
	Password    string `validate:"required" json:"password"`
	Phone       string `validate:"required" json:"phone"`
	Avatar      string `json:"avatar"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
	Created_at  string `json:"created_at"`
}

type UserUpdateBody struct {
	Username    string `validate:"required" json:"username"`
	Fullname    string `validate:"required" json:"fullname"`
	Phone       string `validate:"required" json:"phone"`
	Avatar      string `validate:"required" json:"avatar"`
	Location    string `validate:"required" json:"location"`
	Coordinates string `validate:"required" json:"coordinates"`
}

type User struct {
	Id            int    `json:"id" db:"id"`
	Username      string `json:"username" db:"username"`
	Fullname      string `json:"fullname" db:"fullname"`
	Email         string `json:"email" db:"email"`
	Password      string `json:"password" db:"password"`
	Phone         string `json:"phone" db:"phone"`
	Avatar        string `json:"avatar" db:"avatar"`
	Location      string `json:"location" db:"location"`
	Coordinates   string `json:"coordinates" db:"coordinates"`
	Created_at    string `json:"created_on" db:"created_at"`
	OTP           any    `json:"otp" db:"otp"`
	OTPExpiration any    `json:"otp_expiration" db:"otp_expiration"`
}

type LoginUser struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

type ViewedProductStruct struct {
	Id         int    `json:"id" db:"id"`
	ProductId  string `json:"product_id" db:"product_id"`
	UserId     string `json:"user_id" db:"user_id"`
	Created_at string `json:"created_on" db:"created_at"`
}

type ResetPasswordStruct struct {
	OTP         string `validate:"required" json:"otp"`
	Username    string `validate:"required" json:"username"`
	NewPassword string `validate:"required" json:"new_password"`
}

type ResetRequestStruct struct {
	Username string `validate:"required" json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
