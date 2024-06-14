package product

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/logger"
	validators "olx-clone/functions/validator"
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

	userId, valid := getUserIdFromReq(ctx)
	if !valid {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidUserIdMessage)
	}

	productId := getProductIdFromParams(ctx)
	product, err := model.GetProduct(context.TODO(), productId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.ProductNotFoundMessage)
		}

		logger.WithRequest(ctx).Panicln(err)
	}

	err = model.AddProductViews(context.TODO(), userId, productId)
	if err != nil {
		logger.WithRequest(ctx).Errorln(err)
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

	rows, err := model.GetUsersListPaginatedValue(context.TODO(), itemsPerPage, offset)
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
		"products":       products,
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

	err = validators.ValidateStruct(body)
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
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	productId := ctx.Param("pid")

	if err := model.DeleteProductViewsByProductId(context.TODO(), productId); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err := model.DeleteProductByProductId(context.TODO(), productId); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// get product views
func GetProductViews(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	productId := ctx.Param("pid")
	views, err := model.GetProductViews(context.TODO(), productId);
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.ProductNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"views": views,
	})
}
