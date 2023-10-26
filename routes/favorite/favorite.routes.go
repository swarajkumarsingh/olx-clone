package routes

import (
	"olx-clone/controller/favorite"

	"github.com/gin-gonic/gin"
)

func AddProductRoutes(router *gin.Engine) {
	favorites := router.Group("/")

	favorites.GET("/favorites", favorite.GetAllFavorite)
	favorites.POST("/create/favorite", favorite.AddFavorite)
	favorites.DELETE("/favorite/:fid", favorite.DeleteFavorite)
}
