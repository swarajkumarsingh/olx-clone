package user

import (
	"net/http"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserBody struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Number     string `json:"number"`
	Avatar     string `json:"avatar"`
	Address    string `json:"address"`
	Created_on string `json:"created_on"`
}

func CreateUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, 500)

	var body UserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = InsertUser(ctx, body, hashedPassword); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	tokenString, err := GenerateJwtToken(body.Name)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(200, gin.H{
		"error":   false,
		"message": "User Created successfully",
		"token":   tokenString,
	})
}

// get user
func GetUsers(ctx *gin.Context) {
	errorHandler.Recovery(ctx, http.StatusInternalServerError)

	const pageSize = 10
	database := db.Mgr.DBConn

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(pageSize)))

	// Calculate the offset to determine the starting record for the current page.
	offset := (page - 1) * itemsPerPage

	// Query the users from the database with pagination.
	rows, err := database.Query(`SELECT id, name, email FROM "user" ORDER BY id LIMIT $1 OFFSET $2`, itemsPerPage, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	// Iterate through the rows and build a slice of user data.
	var users []gin.H
	for rows.Next() {
		var id int
		var username, email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user data"})
			return
		}
		users = append(users, gin.H{"id": id, "username": username, "email": email})
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
