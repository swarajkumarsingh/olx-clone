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

	offset := (page - 1) * itemsPerPage
	rows, err := database.Query(`SELECT id, name, email, number FROM "user" ORDER BY id LIMIT $1 OFFSET $2`, itemsPerPage, offset)
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
