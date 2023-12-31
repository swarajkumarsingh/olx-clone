package user

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	model "olx-clone/models/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getCreateUserBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if model.UserAlreadyExistsWithUsername(body.Username) {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.UserAlreadyExistsMessage)
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = model.InsertUser(body, hashedPassword); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	user, err := model.GetUserByUsername(context.TODO(), body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := generateJwtToken(strconv.Itoa(user.Id))
	if err != nil {
		logger.WithRequest(ctx).Panicln("unable to login, try again later")
	}

	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "User Created successfully",
		"token":   token,
	})
}

func GetUsers(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	rows, err := model.GetUsersListPaginatedValue(itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
	}
	defer rows.Close()

	users := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var username, email, number string
		if err := rows.Scan(&id, &username, &email, &number); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
		}
		users = append(users, gin.H{"id": id, "username": username, "email": email, "number": number})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}

// get user
func GetUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username, valid := getUserNameFromParam(ctx)
	if !valid {
		errorHandler.CustomError(ctx, http.StatusBadRequest, messages.InvalidUsernameMessage)
		return
	}

	user, err := model.GetUserByUsername(context.TODO(), username)
	if err == sql.ErrNoRows {
		errorHandler.CustomError(ctx, http.StatusNotFound, messages.UserNotFoundMessage)
		return
	}
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
	})
}

// login
func LoginUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username, password, err := getUserNameAndPasswordFromBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusNotFound, err.Error())
		return
	}

	user, err := model.IsValidUser(context.TODO(), username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := generateJwtToken(strconv.Itoa(user.Id))
	if err != nil {
		logger.WithRequest(ctx).Panicln("unable to login, try again later")
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "User login successfully",
		"token":   token,
	})
}

// test
func Test(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "come home sugar daddy",
	})
}

// update
func UpdateUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getUserUpdateMethodBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	username, err := getUserName(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = model.UpdateUser(context.TODO(), username, body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "User updated successfully",
	})
}

// delete
func DeleteUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username, valid := getUserNameFromParam(ctx)
	if !valid {
		errorHandler.CustomError(ctx, http.StatusBadRequest, messages.InvalidUsernameMessage)
		return
	}

	if err := model.DeleteUserByUsername(username); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// change password
func ResetPasswordDeprecated(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getResetPasswordCredentialsFromBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	err = model.CheckIfCurrentPasswordIsValid(context.TODO(), body.Username, body.NewPassword)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidCredentialsMessage)
	}

	err = model.UpdatePassword(context.TODO(), body.Username, body.NewPassword)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Reset password successful",
	})
}

// request change password
func RequestResetPassword(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getResetRequestCredentialsFromBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	user, err := model.CheckIfUsernameExists(context.TODO(), body.Username)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	otp, err := generateOtp()
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	_, err = sendOtpEmail(constants.DefaultSenderEmailId, user.Email, otp)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	otpSecret := encodeString(otp)
	otpExpiration := getTimeInMinutes(5)
	if err = model.SaveOTPAndExpirationInDB(context.TODO(), user.Username, otpSecret, otpExpiration); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "OTP send successfully",
	})
}

// reset password
func ResetPassword(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getResetPasswordCredentialsFromBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	encodedOTP, err := model.GetOtpFromDB(context.TODO(), body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	otp, err := decodeString(encodedOTP)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	valid := verifyOTPs(body.OTP, otp)
	if !valid {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidOTPMessage)
	}

	err = model.UpdatePassword(context.TODO(), body.Username, body.NewPassword)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	if err = model.ResetOtpAndOtpExpiration(context.TODO(), body.Username); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Reset password successful",
	})
}

// change email
func ChangeEmail(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)
	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "implementation pending",
	})
}

// change phone
func ChangePhoneNumber(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)
	logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, "pending implementation")
}

// view all viewed products
func ViewedProducts(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	userId, valid := getUserIdFromReq(ctx)
	if !valid {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidUserIdMessage)
	}

	rows, err := model.GetViewedProductsListPaginatedValue(userId, itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
	}
	defer rows.Close()

	users := make([]gin.H, 0)
	for rows.Next() {
		var id, userId, productId int
		if err := rows.Scan(&id, &userId, &productId); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
		}
		users = append(users, gin.H{"id": id, "userId": userId, "productId": productId})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}
