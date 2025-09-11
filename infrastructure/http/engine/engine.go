package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/config"
	dep "go.risoftinc.com/xarch/infrastructure/http"
	"go.risoftinc.com/xarch/infrastructure/http/router"
	"gorm.io/gorm"
)

type App struct {
	Config          config.Config
	Logger          gologger.Logger
	DB              *gorm.DB
	ResponseManager *goresponse.ResponseConfig
}

func Start(app App, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Initialize HTTP server
		e := router.Routers(dep.InitializeServices(app.DB, app.Config, app.Logger))

		// Start HTTP server in background
		go func() {
			e.HideBanner = true
			app.Logger.Info(fmt.Sprintf("HTTP server starting on %s:%d", app.Config.Http.Server, app.Config.Http.Port)).Send()
			if err := e.Start(fmt.Sprintf("%s:%d", app.Config.Http.Server, app.Config.Http.Port)); err != nil && err != http.ErrServerClosed {
				app.Logger.Fatal("shutting down the server").Send()
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		app.Logger.Info("Shutdown signal received").Send()

		// Graceful shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown HTTP server
		if err := e.Shutdown(ctx); err != nil {
			app.Logger.Error(fmt.Sprintf("Failed to shutdown HTTP server: %v", err)).Send()
		} else {
			app.Logger.Info("HTTP server shutdown successfully").Send()
		}

		app.Logger.Info("Application shutdown completed").Send()
	}()
}
