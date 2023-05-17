package api

import (
	"context"
	"fmt"
	"github.com/artemys/pprof-visualizer/configs"
	"github.com/artemys/pprof-visualizer/internal/api/router"
	"github.com/artemys/pprof-visualizer/internal/api/services"
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config   configs.ApiConfig
	router   *gin.Engine
	services *services.Services
}

func New() *App {
	app := &App{}
	app.setup()
	return app
}

func (app *App) setup() {
	config := configs.LoadApiConfig()
	svcs := services.InitServices(config)
	r := router.InitializeRouter(svcs, config)
	app.config = config
	app.router = r
	app.services = svcs
}

func (app *App) Run() {
	logger.Log.Info(fmt.Sprintf("RUN APP on PORT %d", app.config.Port))

	// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: app.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Run ListenAndServe", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Run Shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exiting")
}
