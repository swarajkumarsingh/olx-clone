package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"olx-clone/conf"
	"olx-clone/constants"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/user"
	"strconv"
	"time"

	ses "olx-clone/infra/ses"

	sesService "github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Get userId from req.userId
func getUserIdFromReq(ctx *gin.Context) (string, bool) {
	return "1", true
}

func getTimeInMinutes(minute int) time.Time {
	return time.Now().Add(5 * time.Minute)
}

func generateOtp() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = byte(rand.Intn(10)) + 48
	}

	return string(b), nil
}

func encodeString(input string) string {
	data := []byte(input)
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

func decodeString(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func sendOtpEmail(senderId, recipientId, otp string) (*sesService.SendEmailOutput, error) {
	const charSet = "UTF-8"
	const subject = "Forgot password? reset your password using the given OTP."
	var textBody = fmt.Sprintf("Forgot password? <br> <p>Here is your OTP %s from olx, please ignore if you have not requested it", otp)
	var htmlBody = fmt.Sprintf("<h1>Forgot password?</h1> <br> <p>Here is your OTP %s from olx, please ignore if you have not requested it</p>", otp)

	return ses.SendEmail(senderId, recipientId, subject, htmlBody, textBody, charSet)
}

func getCurrentPageValue(ctx *gin.Context) int {
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

func getUserNameFromParam(ctx *gin.Context) (string, bool) {
	username := ctx.Param("username")
	valid := general.ValidUserName(username)
	log.Println(valid)

	if !valid {
		return "", false
	}

	return username, true
}

func getCreateUserBody(ctx *gin.Context) (model.UserBody, error) {
	var body model.UserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return body, err
	}

	if err := general.ValidateStruct(body); err != nil {
		return body, err
	}
	return body, nil
}

func getUserNameAndPasswordFromBody(ctx *gin.Context) (string, string, error) {
	var loginCredentials model.LoginUser
	if err := ctx.ShouldBindJSON(&loginCredentials); err != nil || !general.ValidUserName(loginCredentials.Username) {
		return "", "", errors.New("invalid username or password")
	}
	return loginCredentials.Username, loginCredentials.Password, nil
}

func getResetPasswordCredentialsFromBody(ctx *gin.Context) (model.ResetPasswordStruct, error) {
	var model model.ResetPasswordStruct
	if err := ctx.ShouldBindJSON(&model); err != nil {
		return model, errors.New("invalid body")
	}
	if err := general.ValidateStruct(model); err != nil {
		return model, err
	}
	return model, nil
}

func getResetRequestCredentialsFromBody(ctx *gin.Context) (model.ResetRequestStruct, error) {
	var model model.ResetRequestStruct
	if err := ctx.ShouldBindJSON(&model); err != nil {
		return model, errors.New("invalid body")
	}
	if err := general.ValidateStruct(model); err != nil {
		return model, err
	}
	return model, nil
}

func getOffsetValue(page int, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}

func getItemPerPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(constants.DefaultPerPageSize)))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting item per-page value: ", err)
		return constants.DefaultPerPageSize
	}
	return val
}

func calculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
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
