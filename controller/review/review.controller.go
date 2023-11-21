package review

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	model "olx-clone/models/review"
		validators "olx-clone/functions/validator"


	"github.com/gin-gonic/gin"
)

func CreateReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.CreateReviewStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidBodyMessage)
	}

	if err := validators.ValidateStruct(body); err != nil {
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

	if err := validators.ValidateStruct(body); err != nil {
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

	pid := ctx.Param("pid")
	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	rows, err := model.GetProductReviews(context.TODO(), pid, itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
	}
	defer rows.Close()

	products := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var productId, rating, comment string
		if err := rows.Scan(&id, &productId, &rating, &comment); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
		}
		products = append(products, gin.H{"id": id, "product_id": productId, "rating": rating, "comment": comment})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":    products,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}

func DeleteReview(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	rid := ctx.Param("rid")
	if rid == "" {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidUsernameMessage)
	}

	if err := model.DeleteReview(context.TODO(), rid); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
