package product

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/product"

	"github.com/gin-gonic/gin"
)

// create product
func CreateProduct(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	body, err := getCreateProductBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if err := model.CreateProduct(context.TODO(), body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Product Created successfully",
	})
}

// read product
func GetProduct(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	productId := getProductIdFromParams(ctx)

	product, err := model.GetProduct(context.TODO(), productId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.ProductNotFoundMessage)
		}

		logger.WithRequest(ctx).Panicln(err)
	}

	userId, valid := getUserIdFromReq(ctx)
	if valid {
		err = model.AddProductViews(context.TODO(), userId, productId)
		if err != nil {
			logger.WithRequest(ctx).Errorln(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"product": product,
	})
}

func GetProducts(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	rows, err := model.GetUsersListPaginatedValue(itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveProductsMessage)
	}
	defer rows.Close()

	products := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var title, views, price string
		if err := rows.Scan(&id, &title, &views, &price); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveProductsMessage)
		}
		products = append(products, gin.H{"id": id, "title": title, "views": views, "price": price})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":       products,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}

// update product
func UpdateProduct(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	// get body & validate
	body, err := getUpdateProductBody(ctx)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	err = general.ValidateStruct(body)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	// get product id from params
	pid := ctx.Param("pid")

	// update query
	err = model.UpdateProduct(context.TODO(), pid, body)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	// show response
	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "product updated successfully",
	})
}

// delete product
func DeleteProduct(ctx *gin.Context) {

}
