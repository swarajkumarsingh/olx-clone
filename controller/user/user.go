package user

import (
	"net/http"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	model "olx-clone/models/user"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusInternalServerError)

	var body model.UserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = model.InsertUser(body, hashedPassword); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := GenerateJwtToken(body.Name)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(200, gin.H{
		"error":   false,
		"message": "User Created successfully",
		"token":   token,
	})
}

func GetUsers(ctx *gin.Context) {
	errorHandler.Recovery(ctx, http.StatusInternalServerError)

	page := GetCurrentPageValue(ctx)
	itemsPerPage := GetItemPerPageValue(ctx)
	offset := GetOffsetValue(page, itemsPerPage)

	rows, err := model.GetUsersListPaginatedValue(itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln("Failed to retrieve users")
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id int
		var username, email, number string
		if err := rows.Scan(&id, &username, &email, &number); err != nil {
			logger.WithRequest(ctx).Panicln("Failed to scan user data")
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

}

// login
func LoginUser(ctx *gin.Context) {

}

// logout
func LogoutUser(ctx *gin.Context) {

}

// update
func UpdateUser(ctx *gin.Context) {

}

// delete
func DeleteUser(ctx *gin.Context) {

}
