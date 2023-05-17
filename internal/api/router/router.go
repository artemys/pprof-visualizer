package router

import (
	"github.com/artemys/pprof-visualizer/configs"
	"github.com/artemys/pprof-visualizer/internal/api/handlers"
	"github.com/artemys/pprof-visualizer/internal/api/services"
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	fcmsrouter "github.com/artemys/pprof-visualizer/internal/pkg/logger/gin"
	"time"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(services *services.Services, config configs.ApiConfig) *gin.Engine {
	r := gin.New()

	// Initialize Global Middlewares
	initializeMiddlewares(r)

	// Initialize routes
	initializeRoutes(r, services, config)

	return r
}

func initializeRoutes(r *gin.Engine, services *services.Services, config configs.ApiConfig) {
	r.GET("/health", handlers.Healthcheck())
	r.GET("/", handlers.Index())
	r.POST("/visualize", handlers.Visualize())

	r.Use(fcmsrouter.Ginzap(logger.Log, time.RFC3339, true, false))
	r.NoRoute(handlers.NoRoute)
	r.LoadHTMLGlob("**/templates/*")
}
