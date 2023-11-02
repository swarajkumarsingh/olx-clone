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

	var body model.UserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if model.UserAlreadyExistsWithUsername(body.Username) {
		errorHandler.CustomError(ctx, http.StatusBadRequest, messages.UserAlreadyExistsMessage)
		return
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = model.InsertUser(body, hashedPassword); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := GenerateJwtToken(body.Username)
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

	page := GetCurrentPageValue(ctx)
	itemsPerPage := GetItemPerPageValue(ctx)
	offset := GetOffsetValue(page, itemsPerPage)

	rows, err := model.GetUsersListPaginatedValue(itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln("Failed to retrieve users")
	}
	defer rows.Close()

	users := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var username, email, number string
		if err := rows.Scan(&id, &username, &email, &number); err != nil {
			logger.WithRequest(ctx).Panicln("Failed to retrieve users")
		}
		users = append(users, gin.H{"id": id, "username": username, "email": email, "number": number})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": CalculateTotalPages(page, itemsPerPage),
	})
}

// get user
func GetUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username, valid := GetUserNameFromParam(ctx)
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

	username, password, err := GetUserNameAndPasswordFromBody(ctx)
	if err != nil {
		errorHandler.CustomError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, err = model.IsValidUser(username, password)
	if err == sql.ErrNoRows {
		errorHandler.CustomError(ctx, http.StatusNotFound, messages.UserNotFoundMessage)
		return
	}

	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := GenerateJwtToken(username)

	if err != nil {
		logger.WithRequest(ctx).Panicln("Unable to login, try again later")
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

	username, valid := GetUserNameFromParam(ctx)
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
	logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, "pending implementation")
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
