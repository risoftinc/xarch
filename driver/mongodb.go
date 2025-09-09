package driver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/risoftinc/xarch/config"
)

// buildMongoURI builds MongoDB URI with authentication if provided
func buildMongoURI(cfg config.MongoDBConfig) string {
	if cfg.Username != "" && cfg.Password != "" {
		// If URI already contains authentication, use it as is
		if strings.Contains(cfg.URI, "@") {
			return cfg.URI
		}

		// Build URI with authentication
		// Format: mongodb://username:password@host:port/database
		uri := strings.Replace(cfg.URI, "mongodb://", fmt.Sprintf("mongodb://%s:%s@", cfg.Username, cfg.Password), 1)
		return uri
	}

	return cfg.URI
}

// ConnectMongoDB creates a MongoDB connection
func ConnectMongoDB(cfg config.MongoDBConfig) *mongo.Database {
	defer func() {
		if r := recover(); r != nil {
			log.Panic(fmt.Sprint(r))
		}
	}()

	// Build URI with authentication if provided
	uri := buildMongoURI(cfg)
	log.Printf("Connecting to MongoDB at %s", uri)

	// Set client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxIdleTime)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
	}

	log.Printf("MongoDB connection pool configured: MaxPoolSize=%d, MinPoolSize=%d, MaxIdleTime=%v",
		cfg.MaxPoolSize, cfg.MinPoolSize, cfg.MaxIdleTime)

	return client.Database(cfg.Database)
}

// CloseMongoDB closes the MongoDB connection
func CloseMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Failed to close MongoDB connection: %v", err)
	} else {
		log.Printf("MongoDB connection closed successfully")
	}
}
