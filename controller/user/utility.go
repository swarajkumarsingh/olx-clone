package user

import (
	"errors"
	"log"
	"olx-clone/conf"
	"olx-clone/constants"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetCurrentPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting current page value: ", err)
		return 1
	}
	return val
}

func getUserName(ctx *gin.Context) (string, error) {
	username := ctx.Param("username")
	if username == "" {
		return "", errors.New("invalid username")
	}
	return username, nil
}

func getUserUpdateMethodBody(ctx *gin.Context) (model.UserUpdateBody, error) {
	var body model.UserUpdateBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return body, errors.New("invalid username or password")
	}
	return body, nil
}

func GetUserNameFromParam(ctx *gin.Context) (string, bool) {
	username := ctx.Param("username")
	valid := general.ValidUserName(username)
	log.Println(valid)

	if !valid {
		return "", false
	}

	return username, true
}

func GetUserNameAndPasswordFromBody(ctx *gin.Context) (string, string, error) {
	var loginCredentials model.LoginUser
	if err := ctx.ShouldBindJSON(&loginCredentials); err != nil || !general.ValidUserName(loginCredentials.Username) {
		return "", "", errors.New("invalid username or password")
	}
	return loginCredentials.Username, loginCredentials.Password, nil
}

func GetResetPasswordCredentialsFromBody(ctx *gin.Context) (model.ResetPasswordStruct, error) {
	var model model.ResetPasswordStruct
	if err := ctx.ShouldBindJSON(&model); err != nil {
		return model, errors.New("invalid body")
	}
	if err := general.ValidateStruct(model); err != nil {
		return model, err
	}
	return model, nil
}

func GetOffsetValue(page int, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}

func GetItemPerPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(constants.DefaultPerPageSize)))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting item per-page value: ", err)
		return constants.DefaultPerPageSize
	}
	return val
}

func CalculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptHashingCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwtToken(name string) (string, error) {
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
