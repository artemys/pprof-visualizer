package router

import (
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	fcmsrouter "github.com/artemys/pprof-visualizer/internal/pkg/logger/gin"
	"github.com/gin-gonic/gin"
)

func initializeMiddlewares(router *gin.Engine) {
	// use recovery with zap
	router.Use(fcmsrouter.RecoveryWithZap(logger.Log, true))
}
