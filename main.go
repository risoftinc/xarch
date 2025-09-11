package main

import (
	"log"
	"sync"

	"go.risoftinc.com/gologger"
	"go.risoftinc.com/xarch/config"
	"go.risoftinc.com/xarch/driver"

	grpc "go.risoftinc.com/xarch/infrastructure/grpc/engine"
	http "go.risoftinc.com/xarch/infrastructure/http/engine"
)

func main() {
	// Load configuration
	cfg := config.Configuration()

	// Connect to database using existing driver
	db := driver.ConnectDB(cfg.Database)
	defer driver.CloseDB(db)

	// MongoDB Connection Example (uncomment to use)
	// mongoDB := driver.ConnectMongoDB(cfg.MongoDB)
	// if mongoDB != nil {
	// 	defer driver.CloseMongoDB(mongoDB.Client())
	// 	log.Println("MongoDB connected successfully")
	// }

	// Redis Connection Example (uncomment to use)
	// redisClient := driver.ConnectRedis(cfg.Redis)
	// if redisClient != nil {
	// 	defer driver.CloseRedis(redisClient)
	// 	log.Println("Redis connected successfully")
	// }

	// Load response manager

	// If use sync config manager
	// responseManager, err := driver.ResponseManager(cfg.ResponseManager)
	// if err != nil {
	// 	log.Fatalf("Failed to load response manager: %v", err)
	// }

	responseManager, err := driver.ResponseManagerAsync(cfg.ResponseManager)
	if err != nil {
		log.Fatalf("Failed to load response manager: %v", err)
	}
	defer responseManager.Stop()

	// Initialize logger with config
	logger := gologger.NewLoggerWithConfig(gologger.LoggerConfig{
		OutputMode:   cfg.Logger.OutputMode,
		LogLevel:     cfg.Logger.LogLevel,
		LogDir:       cfg.Logger.LogDir,
		RequestIDKey: "traceID",
		ShowCaller:   true,
	})
	defer logger.Close()

	// Simple approach - just start both servers and wait for signal
	var wg sync.WaitGroup

	// Start HTTP server
	http.Start(http.App{
		Config:          cfg,
		Logger:          logger,
		DB:              db,
		ResponseManager: responseManager.GetConfig(),
	}, &wg)

	// Start GRPC server
	grpc.StartGRPC(grpc.App{
		Config:          cfg,
		Logger:          logger,
		DB:              db,
		ResponseManager: responseManager.GetConfig(),
	}, &wg)

	// Wait for both servers to complete
	wg.Wait()
}
