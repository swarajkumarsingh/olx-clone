package productRoutes

import (
	"olx-clone/authentication"
	"olx-clone/controller/product"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	products := router.Group("/")

	products.POST("/product", product.CreateProduct)
	products.GET("/products",  product.GetProducts)
	products.GET("/product/:pid",authentication.AuthorizeUser, product.GetProduct)
	products.GET("/product/:pid/views", product.GetProductViews)
	products.PATCH("/product/:pid", product.UpdateProduct)
	products.DELETE("/product/:pid", product.DeleteProduct)
}
