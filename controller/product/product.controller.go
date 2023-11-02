package product

import (
	"context"
	"net/http"
	"olx-clone/errorHandler"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/product"

	"github.com/gin-gonic/gin"
)

// create product
func CreateProduct(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.CreateProductBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err.Error())
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

	productId := GetProductIdFromParams(ctx)

	product, err := model.GetProduct(context.TODO(), productId)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	userId, valid := GetUserIdFromReq(ctx)
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

// update product
func UpdateProduct(ctx *gin.Context) {

}

// delete product
func DeleteProduct(ctx *gin.Context) {

}
