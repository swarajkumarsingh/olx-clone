package user

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	model "olx-clone/models/user"

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

	token, err := generateJwtToken(body.Username)
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

	user, err := model.GetUserByUsername(username)
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

	_, err = model.IsValidUser(username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := generateJwtToken(username)

	if err != nil {
		logger.WithRequest(ctx).Panicln("unable to login, try again later")
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "User login successfully",
		"token":   token,
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
func ResetPassword(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getResetPasswordCredentialsFromBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	err = model.CheckIfCurrentPasswordIsValid(context.TODO(), body.Username, body.CurrentPassword)
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

// change email
func ChangeEmail(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)
	logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, "pending implementation")
}

// change phone
func ChangePhoneNumber(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)
	logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, "pending implementation")
}

// view all viewed products
func ViewedProducts(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)
	logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, "pending implementation")
}
