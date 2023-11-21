package favorite

import (
	"context"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	validators "olx-clone/functions/validator"
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

	if err := validators.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if err := model.AddFavorite(context.TODO(), body.UserId, body.ProductId); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "added to favorites successfully",
	})
}

// remove products from favorite
func DeleteFavorite(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	fid := ctx.Param("fid")

	if err := model.RemoveFavorite(context.TODO(), fid); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "removed from favorites successfully",
	})
}

// get all user favorites product
func GetAllUsersFavorite(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	userId, _ := getUserIdFromReq(ctx)
	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	rows, err := model.GetUsersListPaginatedValue(userId, itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
	}
	defer rows.Close()

	users := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var userId, productId string
		if err := rows.Scan(&id, &userId, &productId); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
		}
		users = append(users, gin.H{"id": id, "userId": userId, "productId": productId})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}
