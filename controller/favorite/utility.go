package favorite

import (
	"fmt"
	"olx-clone/constants"
	"olx-clone/functions/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserIdFromReq(ctx *gin.Context) (string, bool) {
	uid, valid := ctx.Get(constants.UserIdMiddlewareConstant)
	if !valid || uid == nil || fmt.Sprintf("%v", uid) == "" {
		return "", false
	}

	return fmt.Sprintf("%v", uid), true
}

func getCurrentPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting current page value: ", err)
		return 1
	}
	return val
}

func getItemPerPageValue(ctx *gin.Context) int {
	val, err := strconv.Atoi(ctx.DefaultQuery("per_page", strconv.Itoa(constants.DefaultPerPageSize)))
	if err != nil {
		logger.WithRequest(ctx).Errorln("error while extracting item per-page value: ", err)
		return constants.DefaultPerPageSize
	}
	return val
}

func getOffsetValue(page int, itemsPerPage int) int {
	return (page - 1) * itemsPerPage
}

func calculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
}
