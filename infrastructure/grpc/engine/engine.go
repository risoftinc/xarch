package engine

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.risoftinc.com/gologger"
	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/config"
	dep "go.risoftinc.com/xarch/infrastructure/grpc"
	"go.risoftinc.com/xarch/infrastructure/grpc/router"
	"gorm.io/gorm"
)

type App struct {
	Config          config.Config
	Logger          gologger.Logger
	DB              *gorm.DB
	ResponseManager *goresponse.ResponseConfig
}

func StartGRPC(app App, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Initialize dependencies
		dependencies := dep.InitializeServices(app.DB, app.Config, app.Logger)

		// Register services
		grpcServer := router.RegisterGRPCServices(dependencies)

		// Start gRPC server in background
		go func() {
			// Create listener
			lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.Config.Grpc.Server, app.Config.Grpc.Port))
			if err != nil {
				app.Logger.Fatal("Failed to listen on gRPC port: " + err.Error()).Send()
			}

			app.Logger.Info(fmt.Sprintf("gRPC server starting on %s:%d", app.Config.Grpc.Server, app.Config.Grpc.Port)).Send()

			if err := grpcServer.Serve(lis); err != nil {
				app.Logger.Fatal("Failed to serve gRPC server: " + err.Error()).Send()
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

		// Shutdown gRPC server
		done := make(chan bool)
		go func() {
			grpcServer.GracefulStop()
			done <- true
		}()

		select {
		case <-done:
			app.Logger.Info("gRPC server shutdown successfully").Send()
		case <-ctx.Done():
			app.Logger.Info("gRPC server shutdown timeout, forcing stop").Send()
			grpcServer.Stop()
		}

		app.Logger.Info("gRPC Application shutdown completed").Send()
	}()
}
