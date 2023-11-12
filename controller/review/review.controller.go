package review

import (
	"context"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/review"

	"github.com/gin-gonic/gin"
)

func CreateReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.CreateReviewStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if err := model.AddReview(context.TODO(), body.UserId, body.ProductId, body.Rating, body.Comment); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "review added successfully",
	})
}

func GetReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"review": "",
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
