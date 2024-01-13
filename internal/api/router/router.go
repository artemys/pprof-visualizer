package router

import (
	"github.com/artemys/pprof-visualizer/internal/api/handlers"
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	fcmsrouter "github.com/artemys/pprof-visualizer/internal/pkg/logger/gin"
	"github.com/gin-contrib/pprof"
	"time"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	r := gin.New()

	// Initialize Global Middlewares
	initializeMiddlewares(r)

	// Initialize routes
	initializeRoutes(r)

	return r
}

func initializeRoutes(r *gin.Engine) {
	r.GET("/health", handlers.Healthcheck())
	r.GET("/", handlers.Index())
	r.POST("/visualize", handlers.Visualize())

	r.Use(fcmsrouter.Ginzap(logger.Log, time.RFC3339, true, false))
	r.NoRoute(handlers.NoRoute)
	r.LoadHTMLGlob("**/templates/*")

	pprof.Register(r)
}
