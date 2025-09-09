package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/risoftinc/gologger"
	"github.com/risoftinc/goresponse"
	"github.com/risoftinc/xarch/config"
	dep "github.com/risoftinc/xarch/infrastructure/http"
	"github.com/risoftinc/xarch/infrastructure/http/router"
	"gorm.io/gorm"
)

type App struct {
	Config          config.Config
	Logger          gologger.Logger
	DB              *gorm.DB
	ResponseManager *goresponse.ResponseConfig
}

func Start(app App) {
	// Initialize HTTP server
	e := router.Routers(dep.InitializeServices(app.DB, app.Config, app.Logger))

	// Start HTTP server in background
	go func() {
		e.HideBanner = true
		if err := e.Start(fmt.Sprintf("%s:%d", app.Config.Http.Server, app.Config.Http.Port)); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutdown signal received")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := e.Shutdown(ctx); err != nil {
		app.Logger.Error(fmt.Sprintf("Failed to shutdown HTTP server: %v", err))
	} else {
		app.Logger.Info("HTTP server shutdown successfully")
	}

	app.Logger.Info("Application shutdown completed")
}
