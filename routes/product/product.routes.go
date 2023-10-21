package routes

import (
	"github.com/gin-gonic/gin"
)

func AddProductRoutes(router *gin.Engine) {
	users := router.Group("/")

	users.POST("/create-product", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"error": false,
		})
	})
}
