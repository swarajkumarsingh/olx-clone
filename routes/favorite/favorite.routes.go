package favoriteRoutes

import (
	"olx-clone/controller/favorite"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	favorites := router.Group("/")

	favorites.GET("/favorites", favorite.GetAllUsersFavorite)
	favorites.POST("/favorite", favorite.AddFavorite)
	favorites.DELETE("/favorite/:fid", favorite.DeleteFavorite)
}
