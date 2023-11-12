package reviewRoutes

import (
	"olx-clone/controller/favorite"
	"olx-clone/controller/review"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	reviews := router.Group("/")

	reviews.POST("/review", review.CreateReview)
	reviews.GET("/review/:rid", review.GetReview)
	reviews.PATCH("/review", review.UpdateReview)
	reviews.DELETE("/review/:rid", favorite.DeleteFavorite)
	reviews.GET("/product/review/:pid", favorite.GetAllUsersFavorite)
}
