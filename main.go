package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"olx-clone/functions/logger"
	userRoutes "olx-clone/routes/user"
	reviewRoutes "olx-clone/routes/review"
	sellerRoutes "olx-clone/routes/seller"
	productRoutes "olx-clone/routes/product"
	favoriteRoutes "olx-clone/routes/favorite"
)

var log = logger.Log
var version string = "1.0"

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "app_total_requests",
			Help: "Total number of requests to my app",
		},
	)
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius_a",
	Help: "Current temperature of the CPU.",
})

func CustomMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cpuTemp.Set(float64(100))
		totalRequests.Inc()
		c.Next()
	}
}

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(totalRequests)
}

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
	r.Use(CustomMetricsMiddleware())

	// run migrations
	MigrateDB()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"error":   false,
			"message": "health ok",
		})
	})
	r.GET("/metrics", prometheusHandler())

	userRoutes.AddRoutes(r)
	sellerRoutes.AddRoutes(r)
	reviewRoutes.AddRoutes(r)
	productRoutes.AddRoutes(r)
	favoriteRoutes.AddRoutes(r)

	log.Printf("Server Started, version: %s", version)
	http.ListenAndServe(":8080", r)
}
