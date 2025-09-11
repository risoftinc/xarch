package config

import (
	"fmt"
	"log"
	"time"

	env "go.risoftinc.com/goenv"
)

type (
	Config struct {
		Http            HttpServer
		Grpc            GrpcServer
		Database        DatabaseConfig
		MongoDB         MongoDBConfig
		Redis           RedisConfig
		Logger          LoggerConfig
		ResponseManager ResponseManager
	}

	HttpServer struct {
		Server string
		Port   int
		URL    string
	}

	GrpcServer struct {
		Server string
		Port   int
		URL    string
	}

	DatabaseConfig struct {
		Type          string // "postgres", "mysql", "sqlite"
		PostgresDB    PostgresDB
		MySQLDB       MySQLDB
		SQLiteDB      SQLiteDB
		DBMaxIdleCon  int
		DBMaxOpenCon  int
		DBMaxLifeTime int
		DBDebug       bool
	}

	PostgresDB struct {
		DBUser     string
		DBPass     string
		DBServer   string
		DBPort     int
		DBName     string
		DBTimeZone string
		SSLMode    string
	}

	MySQLDB struct {
		DBUser     string
		DBPass     string
		DBServer   string
		DBPort     int
		DBName     string
		DBTimeZone string
		Charset    string
		ParseTime  bool
		Loc        string
	}

	SQLiteDB struct {
		DBPath string
	}

	MongoDBConfig struct {
		URI         string
		Username    string
		Password    string
		Database    string
		MaxPoolSize uint64
		MinPoolSize uint64
		MaxIdleTime time.Duration
		Timeout     time.Duration
	}

	RedisConfig struct {
		Host         string
		Port         int
		Username     string
		Password     string
		DB           int
		MaxRetries   int
		PoolSize     int
		MinIdleConns int
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	LoggerConfig struct {
		OutputMode string
		LogLevel   string
		LogDir     string
	}

	ResponseManager struct {
		Method   string
		Path     string
		Interval time.Duration
	}
)

func Configuration() Config {
	if err := env.LoadEnv(".env"); err != nil {
		log.Println("error read .env file %w", err.Error())
	}

	cfg := Config{
		Http:            loadHttpServer(),
		Grpc:            loadGrpcServer(),
		Database:        loadDatabaseConfig(),
		MongoDB:         loadMongoDBConfig(),
		Redis:           loadRedisConfig(),
		Logger:          loadLoggerConfig(),
		ResponseManager: loadResponseManagerConfig(),
	}

	log.Println("Success for load all configuration")

	return cfg
}

func loadHttpServer() HttpServer {
	var cfg HttpServer

	cfg.Server = env.GetEnv("SERVER", "localhost")
	cfg.Port = env.GetEnv("PORT", 9000)
	if env.GetEnv("USING_SECURE", true) {
		cfg.URL = "https://" + cfg.Server
	} else {
		cfg.URL = "http://" + cfg.Server
	}

	if cfg.Port != 0 {
		cfg.URL += fmt.Sprintf(":%d", cfg.Port)
	}
	cfg.URL += "/"

	return cfg
}

func loadGrpcServer() GrpcServer {
	var cfg GrpcServer

	cfg.Server = env.GetEnv("GRPC_SERVER", "localhost")
	cfg.Port = env.GetEnv("GRPC_PORT", 9001)
	cfg.URL = fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)

	return cfg
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Type: env.GetEnv("DB_TYPE", "postgres"),
		PostgresDB: PostgresDB{
			DBUser:     env.GetEnv("DB_USER", "root"),
			DBPass:     env.GetEnv("DB_PASS", ""),
			DBServer:   env.GetEnv("DB_SERVER", "localhost"),
			DBPort:     env.GetEnv("DB_PORT", 5432),
			DBName:     env.GetEnv("DB_NAME", "public"),
			DBTimeZone: env.GetEnv("DB_TIME_ZONE", "Asia/Jakarta"),
			SSLMode:    env.GetEnv("DB_SSL_MODE", "disable"),
		},
		MySQLDB: MySQLDB{
			DBUser:     env.GetEnv("DB_USER", "root"),
			DBPass:     env.GetEnv("DB_PASS", ""),
			DBServer:   env.GetEnv("DB_SERVER", "localhost"),
			DBPort:     env.GetEnv("DB_PORT", 3306),
			DBName:     env.GetEnv("DB_NAME", "public"),
			DBTimeZone: env.GetEnv("DB_TIME_ZONE", "Asia/Jakarta"),
			Charset:    env.GetEnv("DB_CHARSET", "utf8mb4"),
			ParseTime:  env.GetEnv("DB_PARSE_TIME", true),
			Loc:        env.GetEnv("DB_LOC", "Local"),
		},
		SQLiteDB: SQLiteDB{
			DBPath: env.GetEnv("DB_PATH", "database.db"),
		},
		DBMaxIdleCon:  env.GetEnv("DB_MAX_IDLE_CON", 10),
		DBMaxOpenCon:  env.GetEnv("DB_MAX_OPEN_CON", 100),
		DBMaxLifeTime: env.GetEnv("DB_MAX_LIFE_TIME", 10),
		DBDebug:       env.GetEnv("DB_DEBUG", false),
	}
}

func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		OutputMode: env.GetEnv("LOG_OUTPUT_MODE", "both"), // "terminal", "file", "both"
		LogLevel:   env.GetEnv("LOG_LEVEL", "debug"),      // "debug", "info", "warn", "error"
		LogDir:     env.GetEnv("LOG_DIR", "logger"),       // directory for log files
	}
}
func loadMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		URI:         env.GetEnv("MONGODB_URI", "mongodb://localhost:27017"),
		Username:    env.GetEnv("MONGODB_USERNAME", ""),
		Password:    env.GetEnv("MONGODB_PASSWORD", ""),
		Database:    env.GetEnv("MONGODB_DATABASE", "xarch"),
		MaxPoolSize: env.GetEnv("MONGODB_MAX_POOL_SIZE", uint64(100)),
		MinPoolSize: env.GetEnv("MONGODB_MIN_POOL_SIZE", uint64(5)),
		MaxIdleTime: env.GetEnv("MONGODB_MAX_IDLE_TIME", 30*time.Minute),
		Timeout:     env.GetEnv("MONGODB_TIMEOUT", 10*time.Second),
	}
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         env.GetEnv("REDIS_HOST", "localhost"),
		Port:         env.GetEnv("REDIS_PORT", 6379),
		Username:     env.GetEnv("REDIS_USERNAME", "root"),
		Password:     env.GetEnv("REDIS_PASSWORD", ""),
		DB:           env.GetEnv("REDIS_DB", 0),
		MaxRetries:   env.GetEnv("REDIS_MAX_RETRIES", 3),
		PoolSize:     env.GetEnv("REDIS_POOL_SIZE", 10),
		MinIdleConns: env.GetEnv("REDIS_MIN_IDLE_CONNS", 5),
		DialTimeout:  env.GetEnv("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  env.GetEnv("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: env.GetEnv("REDIS_WRITE_TIMEOUT", 3*time.Second),
		IdleTimeout:  env.GetEnv("REDIS_IDLE_TIMEOUT", 5*time.Minute),
	}
}

func loadResponseManagerConfig() ResponseManager {
	return ResponseManager{
		Method:   env.GetEnv("RESPONSE_MANAGER_METHOD", "file"),             // "file", "http"
		Path:     env.GetEnv("RESPONSE_MANAGER_PATH", "config/config.json"), // path to the response manager config file
		Interval: env.GetEnv("RESPONSE_MANAGER_INTERVAL", 5*time.Minute),    // interval in duration
	}
}
