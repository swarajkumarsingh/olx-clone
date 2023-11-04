package product

import (
	"errors"
	"olx-clone/constants"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProductIdFromParams(ctx *gin.Context) string {
	return ctx.Param("pid")
}

func getCreateProductBody(ctx *gin.Context) (model.CreateProductBody, error) {
	var body model.CreateProductBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return body, errors.New("invalid body")
	}

	if err := general.ValidateStruct(body); err != nil {
		return body, err
	}

	return body, nil
}

// TODO: Get userId from req.userId
func getUserIdFromReq(ctx *gin.Context) (string, bool) {
	return "1", true
}

func getCurrentPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting current page value: ", err)
		return 1
	}
	return val
}

func getOffsetValue(page int, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}

func getItemPerPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(constants.DefaultPerPageSize)))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting item per-page value: ", err)
		return constants.DefaultPerPageSize
	}
	return val
}

func calculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
}
