package userRoutes

import (
	"olx-clone/controller/user"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	users := router.Group("/")

	users.POST("/login", user.LoginUser)
	
	users.POST("/user", user.CreateUser)
	users.GET("/users", user.GetUsers)
	users.GET("/user/:username", user.GetUser)
	users.PATCH("/user/:username", user.UpdateUser)
	users.DELETE("/user/:username", user.DeleteUser)
}
