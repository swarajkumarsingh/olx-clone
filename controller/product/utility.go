package product

import "github.com/gin-gonic/gin"

func GetProductIdFromParams(ctx *gin.Context) string {
	return ctx.Param("pid")
}

// TODO: Get userId from req.userId
func GetUserIdFromReq(ctx *gin.Context) (string, bool) {
	return "1", true
}