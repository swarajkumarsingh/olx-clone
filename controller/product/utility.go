package product

import (
	"olx-clone/constants"
	"olx-clone/functions/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProductIdFromParams(ctx *gin.Context) string {
	return ctx.Param("pid")
}

// TODO: Get userId from req.userId
func GetUserIdFromReq(ctx *gin.Context) (string, bool) {
	return "1", true
}

func GetCurrentPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting current page value: ", err)
		return 1
	}
	return val
}


func GetOffsetValue(page int, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}


func GetItemPerPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(constants.DefaultPerPageSize)))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting item per-page value: ", err)
		return constants.DefaultPerPageSize
	}
	return val
}

func CalculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
}