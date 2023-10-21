package user

import (
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"

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

// create user
func CreateUser(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, 500)

	// get body data
	var body UserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	// hash password
	

	// store in DB
	database := db.Mgr.DBConn
	query := `INSERT INTO "user"(name, email, password, number) VALUES($1, $2, $3, $4)`
	result, err := database.Exec(query, body.Name, body.Email, body.Password, body.Number)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}
	logger.Log.Println(result.LastInsertId())
	logger.Log.Println(result.RowsAffected())

	// generate JWT token

	ctx.JSON(200, gin.H{
		"error":   false,
		"message": "User Created successfully",
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
