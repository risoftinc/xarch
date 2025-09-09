package main

import (
	"log"

	"github.com/risoftinc/gologger"
	"github.com/risoftinc/xarch/config"
	"github.com/risoftinc/xarch/driver"
	http "github.com/risoftinc/xarch/infrastructure/http/engine"
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
	log := gologger.NewLoggerWithConfig(gologger.LoggerConfig{
		OutputMode:   cfg.Logger.OutputMode,
		LogLevel:     cfg.Logger.LogLevel,
		LogDir:       cfg.Logger.LogDir,
		RequestIDKey: "traceID",
	})
	defer log.Close()

	// Start HTTP server with graceful shutdown
	http.Start(http.App{
		Config: cfg,
		Logger: log,
		DB:     db,
		// ResponseManager: responseManager, -> If use sync config manager
		ResponseManager: responseManager.GetConfig(),
	})
}
