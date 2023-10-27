package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"olx-clone/functions/logger"
	favoriteRoutes "olx-clone/routes/favorite"
	productRoutes "olx-clone/routes/product"
	userRoutes "olx-clone/routes/user"
)

var log = logger.Log
var version string = "1.0"

func enableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Api-Key, token, User-Agent, Referer")
		c.Writer.Header().Set("AllowCredentials", "true")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		if c.Request.Method == "OPTIONS" {
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()	

	// custom middleware
	r.Use(enableCORS())

	// run migrations
	MigrateDB()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health ok",
		})
	})

	userRoutes.AddRoutes(r)
	productRoutes.AddRoutes(r)
	favoriteRoutes.AddRoutes(r)

	log.Printf("Server Started, version: %s", version)
	http.ListenAndServe(":8080", r)
}
