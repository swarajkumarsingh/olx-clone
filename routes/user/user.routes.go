package userRoutes

import (
	"olx-clone/authentication"
	"olx-clone/controller/user"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	users := router.Group("/")

	users.GET("/test", authentication.AuthorizeUser, user.Test)
	users.POST("/login", user.LoginUser)

	users.POST("/request/password", user.RequestResetPassword)
	users.POST("/verify/password", user.ResetPassword)

	users.POST("/reset/email", user.ChangeEmail)

	users.POST("/user", user.CreateUser)
	users.GET("/users", user.GetUsers)
	users.GET("/user/:username", user.GetUser)
	users.PATCH("/user/:username", user.UpdateUser)
	users.DELETE("/user/:username", user.DeleteUser)

	users.GET("/user/viewed", authentication.AuthorizeUser, user.ViewedProducts)
}
