package productRoutes

import (
	"olx-clone/controller/product"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	products := router.Group("/")

	products.POST("/product/create", product.CreateProduct)
	products.GET("/product/:pid", product.GetProduct)
	products.PATCH("/product/:pid", product.UpdateProduct)
	products.DELETE("/product/:pid", product.DeleteProduct)
}
