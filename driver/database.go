package driver

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.risoftinc.com/xarch/config"
)

// ConnectDB creates a database connection based on the database type
func ConnectDB(cfg config.DatabaseConfig) *gorm.DB {
	defer func() {
		if r := recover(); r != nil {
			log.Panic(fmt.Sprint(r))
		}
	}()

	log.Printf("Connecting to %s database", cfg.Type)

	var db *gorm.DB
	var err error

	switch cfg.Type {
	case "postgres":
		db, err = connectPostgres(cfg.PostgresDB)
	case "mysql":
		db, err = connectMySQL(cfg.MySQLDB)
	case "sqlite":
		db, err = connectSQLite(cfg.SQLiteDB)
	default:
		panic(fmt.Sprintf("Unsupported database type: %s", cfg.Type))
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to %s database: %v", cfg.Type, err))
	}

	if cfg.DBDebug {
		db = db.Debug()
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get underlying sql.DB")
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleCon)                                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenCon)                                  // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifeTime) * time.Minute) // Maximum lifetime of connections

	log.Printf("Database connection pool configured: MaxIdle=%d, MaxOpen=%d, MaxLifetime=%dmin",
		cfg.DBMaxIdleCon, cfg.DBMaxOpenCon, cfg.DBMaxLifeTime)

	return db
}

// connectPostgres creates a PostgreSQL connection
func connectPostgres(cfg config.PostgresDB) (*gorm.DB, error) {
	// Determine SSL mode string
	sslMode := "disable"
	if cfg.SSLMode != "" {
		sslMode = cfg.SSLMode
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBServer,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
		sslMode,
		cfg.DBTimeZone,
	)

	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
}

// connectMySQL creates a MySQL connection
func connectMySQL(cfg config.MySQLDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBServer,
		cfg.DBPort,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// connectSQLite creates a SQLite connection
func connectSQLite(cfg config.SQLiteDB) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
}

// CloseDBWithTimeout closes the database connection with a timeout
// This function should be used in a defer statement for graceful shutdown
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sqlDB, err := db.WithContext(ctx).DB()
	if err != nil {
		log.Printf("Failed to get DB connection: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close DB connection: %v", err)
	} else {
		log.Printf("Database connection closed successfully")
	}
}
