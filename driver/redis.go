package driver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"go.risoftinc.com/xarch/config"
)

// ConnectRedis creates a Redis connection
func ConnectRedis(cfg config.RedisConfig) *redis.Client {
	defer func() {
		if r := recover(); r != nil {
			log.Panic(fmt.Sprint(r))
		}
	}()

	log.Printf("Connecting to Redis at %s:%d", cfg.Host, cfg.Port)

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	log.Printf("Redis connection pool configured: PoolSize=%d, MinIdleConns=%d, MaxRetries=%d",
		cfg.PoolSize, cfg.MinIdleConns, cfg.MaxRetries)

	return rdb
}

// CloseRedis closes the Redis connection
func CloseRedis(client *redis.Client) {
	if client == nil {
		return
	}

	if err := client.WithTimeout(30 * time.Second).Close(); err != nil {
		log.Printf("Failed to close Redis connection: %v", err)
	} else {
		log.Printf("Redis connection closed successfully")
	}
}
