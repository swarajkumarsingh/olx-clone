package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"olx-clone/conf"
	"olx-clone/functions/logger"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	sentrygin "github.com/getsentry/sentry-go/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	ddprofiler "gopkg.in/DataDog/dd-trace-go.v1/profiler"
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
	// init sentry
	if conf.ENV == conf.ENV_PROD || strings.HasPrefix(conf.ENV, "uat") {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         conf.SentryDSN,
			Environment: conf.ENV,
		})
		if err != nil {
			log.Panicf("sentry.Init: %s", err)
		}
		defer sentry.Flush(time.Second)
	}

	// add middlewares
	sentryMiddleware := sentrygin.New(
		sentrygin.Options{
			Repanic: true, // this is important if Recoverer middleware is used
		},
	)

	// do not start datadog in local
	if conf.ENV != conf.ENV_LOCAL {
		// starting datadog tracer
		ddtracer.Start(
			ddtracer.WithAgentAddr(fmt.Sprintf("%s:8126", conf.DDAgentHost)),
			ddtracer.WithDogstatsdAddress(fmt.Sprintf("%s:8125", conf.DDAgentHost)),
			ddtracer.WithEnv(conf.ClientENV),
			ddtracer.WithRuntimeMetrics(),
			ddtracer.WithService(conf.DDServiceName),
		)
		defer ddtracer.Stop()

		// starting datadog profiler
		err := ddprofiler.Start(
			ddprofiler.WithAgentAddr(fmt.Sprintf("%s:8126", conf.DDAgentHost)),
			ddprofiler.WithEnv(conf.ClientENV),
			ddprofiler.WithService(conf.DDServiceName),
			ddprofiler.WithProfileTypes(
				ddprofiler.CPUProfile,
				ddprofiler.HeapProfile,
				ddprofiler.BlockProfile,
				ddprofiler.GoroutineProfile,
				ddprofiler.MetricsProfile,
				ddprofiler.MutexProfile,
			),
		)
		if err != nil {
			log.Panicf("error while starting datadog profiler: %s", err)
		}
		defer ddprofiler.Stop()
	}

	r := gin.Default()

	// skipping datadog middleware in local
	if conf.ENV != conf.ENV_LOCAL {
		r.Use(gintrace.Middleware(conf.DDServiceName, gintrace.WithAnalytics(true)))
	}

	r.Use(sentryMiddleware)

	// custom middleware
	r.Use(enableCORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health ok",
		})
	})

	log.Printf("Server Started, version: %s", version)
	http.ListenAndServe(":8080", r)
}
