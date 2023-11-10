package seller

import (
	"errors"
	"olx-clone/conf"
	"olx-clone/constants"
	"olx-clone/constants/messages"
	"olx-clone/functions/general"
	model "olx-clone/models/seller"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func getCreateSellerBody(ctx *gin.Context) (model.SellerBody, error) {
	var body model.SellerBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return body, errors.New(messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		return body, err
	}
	return body, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptHashingCost)
	return string(bytes), err
}

func generateJwtToken(name string) (string, error) {
	expirationTime := time.Now().Add(5 * 24 * time.Hour)
	claims := &model.Claims{
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(conf.JWTSecretKey)
	if err != nil {
		return "", err
	}

	if tokenString == "" {
		return "", errors.New("error while authorizing")
	}

	return tokenString, nil
}
