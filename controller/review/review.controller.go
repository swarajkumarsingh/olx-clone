package review

import (
	"context"
	"database/sql"
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

	rid := ctx.Param("rid")

	review, err := model.GetReview(context.TODO(), rid)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.ReviewNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":  false,
		"review": review,
	})
}

func UpdateReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.UpdateReviewStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if err := model.UpdateReview(context.TODO(), body.ReviewId, body.Rating, body.Comment); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "updated favorites successfully",
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
