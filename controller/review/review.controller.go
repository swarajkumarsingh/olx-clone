package review

import (
	"net/http"
	"olx-clone/errorHandler"

	"github.com/gin-gonic/gin"
)

func CreateReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

func GetReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

func UpdateReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

func GetProductReviews(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

func DeleteReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}
