package userRoutes

import (
	"olx-clone/controller/user"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	users := router.Group("/")

	users.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"error": false,
		})
	})

	users.POST("/create/user", user.CreateUser)
	users.GET("/users", user.GetUsers)
	users.POST("/user/:username", user.UpdateUser)
	users.GET("/user/:username", user.GetUser)
	users.POST("/login", user.LoginUser)
	users.POST("/logout", user.LogoutUser)
	users.POST("/delete-user", user.DeleteUser)
}
