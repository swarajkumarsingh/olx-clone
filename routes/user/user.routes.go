package routes

import (
	"olx-clone/controller/user"

	"github.com/gin-gonic/gin"
)

func AddUserRoutes(router *gin.Engine) {
	users := router.Group("/")

	users.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"error": false,
		})
	})

	users.POST("/create-user", user.CreateUser)
	users.GET("/users", user.GetUsers)
	users.POST("/user/:username", user.CreateUser)
	users.POST("/login", user.CreateUser)
	users.POST("/logout", user.CreateUser)
	users.POST("/delete-user", user.CreateUser)
}
