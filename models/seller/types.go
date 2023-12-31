package sellerModel

import "github.com/golang-jwt/jwt/v5"

type SellerBody struct {
	Avatar        string `json:"avatar"`
	Username      string `validate:"required" json:"username"`
	Fullname      string `validate:"required" json:"fullname"`
	Description   string `validate:"required" json:"description"`
	Email         string `validate:"required" json:"email"`
	Password      string `validate:"required" json:"password"`
	Phone         string `validate:"required" json:"phone"`
	Location      string `validate:"required" json:"location"`
	Coordinates   string `validate:"required" json:"coordinates"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	Rating        string `json:"rating"`
	AccountStatus string `json:"account_status"`
}

type Seller struct {
	Id            int    `json:"id" db:"id"`
	Username      string `json:"username" db:"username"`
	Description   string `json:"description" db:"description"`
	Fullname      string `json:"fullname" db:"fullname"`
	Email         string `json:"email" db:"email"`
	Password      string `json:"password" db:"password"`
	Phone         string `json:"phone" db:"phone"`
	Avatar        string `json:"avatar" db:"avatar"`
	Location      string `json:"location" db:"location"`
	Coordinates   string `json:"coordinates" db:"coordinates"`
	City          any    `json:"city" db:"city"`
	State         any    `json:"state" db:"state"`
	Country       any    `json:"country" db:"country"`
	ZipCode       any    `json:"zip_code" db:"zip_code"`
	Verified      string `json:"is_verified" db:"is_verified"`
	Rating        string `json:"rating" db:"rating"`
	AccountStatus string `json:"account_status" db:"account_status"`
	Created_at    string `json:"created_on" db:"created_at"`
	OTP           any    `json:"otp" db:"otp"`
	OTPExpiration any    `json:"otp_expiration" db:"otp_expiration"`
}

type SellerUpdateStruct struct {
	Username    string `json:"username" db:"username"`
	Description string `json:"description" db:"description"`
	Fullname    string `json:"fullname" db:"fullname"`
	Phone       string `json:"phone" db:"phone"`
	Avatar      string `json:"avatar" db:"avatar"`
	Location    string `json:"location" db:"location"`
	Coordinates string `json:"coordinates" db:"coordinates"`
	City        any    `json:"city" db:"city"`
	State       any    `json:"state" db:"state"`
	Country     any    `json:"country" db:"country"`
	ZipCode     any    `json:"zip_code" db:"zip_code"`
}

type ResetRequestStruct struct {
	Username string `validate:"required" json:"username"`
}

type UpdateAccountStatusStruct struct {
	Username string `validate:"required" json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginUser struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

type ResetPasswordStruct struct {
	OTP         string `validate:"required" json:"otp"`
	Username    string `validate:"required" json:"username"`
	NewPassword string `validate:"required" json:"new_password"`
}

type ReportAccountStruct struct {
	Username string `validate:"required" json:"username"`
	UserId   string `validate:"required" json:"user_id"`
	SellerId string `validate:"required" json:"seller_id"`
	Message  string `validate:"required" json:"message"`
}
