package sellerRoutes

import (
	"olx-clone/controller/seller"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	sellers := router.Group("/")

	sellers.POST("/seller/login", seller.LoginSeller)
	sellers.POST("/seller/report", seller.ReportSeller)

	sellers.PATCH("/seller/ban", seller.BanSeller)
	sellers.PATCH("/seller/suspend", seller.SuspendSeller)
	sellers.PATCH("/seller/activate", seller.ActivateSeller)

	sellers.PATCH("/seller/verify", seller.VerifySeller)

	sellers.GET("/seller/products/:sid", seller.GetAllCreatedProduct)
	
	sellers.POST("/seller", seller.CreateSeller)
	sellers.GET("/sellers", seller.GetAllSeller)
	sellers.GET("/seller/:sid", seller.GetSeller)
	sellers.PATCH("/seller/:sid", seller.UpdateSeller)
	sellers.DELETE("/seller/:sid", seller.DeleteSeller)

	sellers.POST("/seller/verify/password", seller.ResetPasswordSeller)
	sellers.POST("/seller/request/password", seller.RequestResetPassword)
}
