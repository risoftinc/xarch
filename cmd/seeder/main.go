package main

import (
	"log"

	"github.com/risoftinc/goseeder"
	"github.com/risoftinc/xarch/config"
	"github.com/risoftinc/xarch/database/seeders"
	"github.com/risoftinc/xarch/driver"
)

func main() {
	// Load configuration
	cfg := config.Configuration()

	// Connect to database using existing driver
	db := driver.ConnectDB(cfg.Database)

	// Create seeder manager
	manager := goseeder.NewSeederManager()

	// Create main seeder and register all seeders
	mainSeeder := seeders.NewMainSeeder(db)
	mainSeeder.RegisterAll(manager)

	// Create CLI and run
	cli := goseeder.NewCLIWithAppName(manager, "go run cmd/seeder/main.go")
	if err := cli.Run(); err != nil {
		log.Fatalf("Seeder error: %v", err)
	}

	log.Println("Seeder completed successfully!")
}
