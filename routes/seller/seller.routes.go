package sellerRoutes

import (
	"olx-clone/controller/seller"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	sellers := router.Group("/")

	sellers.POST("/seller/login", seller.LoginSeller)

	sellers.POST("/suspend/seller/:sid", seller.SuspendSeller)
	sellers.POST("/unsuspend/seller/:sid", seller.UnSuspendSeller)

	sellers.POST("/ban/seller/:sid", seller.BanSeller)
	sellers.POST("/unban/seller/:sid", seller.UnBanSeller)

	sellers.POST("/seller/verify", seller.VerifySeller)

	sellers.GET("/seller/products", seller.GetAllCreatedProduct)

	sellers.POST("/seller", seller.CreateSeller)
	sellers.GET("/sellers", seller.GetAllSeller)
	sellers.GET("/seller/:sid", seller.GetSeller)
	sellers.PATCH("/seller/:sid", seller.UpdateSeller)
	sellers.DELETE("/seller/:sid", seller.DeleteSeller)

	sellers.POST("/seller/verify/password", seller.ResetPasswordSeller)
	sellers.POST("/seller/request/password", seller.RequestResetPassword)
}
