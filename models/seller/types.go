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

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
