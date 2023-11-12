package reviewRoutes

import (
	"olx-clone/controller/review"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	reviews := router.Group("/")

	reviews.POST("/review", review.CreateReview)
	reviews.GET("/review/:rid", review.GetReview)
	reviews.PATCH("/review", review.UpdateReview)
	reviews.DELETE("/review/:rid", review.DeleteReview)
	reviews.GET("/product/review/:pid", review.GetProductReviews)
}
