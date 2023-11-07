package favorite

import (
	"context"
	"net/http"
	"olx-clone/errorHandler"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/favorite"

	"github.com/gin-gonic/gin"
)

// add product to favorite
func AddFavorite(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.CreateUserStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	err := model.AddFavorite(context.TODO(), body.UserId, body.ProductId)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

// remove products from favorite
func DeleteFavorite(ctx *gin.Context) {

}

// get all user favorites product
func GetAllFavorite(ctx *gin.Context) {

}
