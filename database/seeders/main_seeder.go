package seeders

import (
	"log"

	"go.risoftinc.com/goseeder"
	"gorm.io/gorm"
)

// MainSeeder contains all seeder methods
type MainSeeder struct {
	db *gorm.DB
}

// NewMainSeeder creates a new main seeder instance
func NewMainSeeder(db *gorm.DB) *MainSeeder {
	return &MainSeeder{db: db}
}

// RegisterAll registers all available seeders with the given manager
func (s *MainSeeder) RegisterAll(manager *goseeder.SeederManager) {
	log.Println("Registering all seeders...")

	// Register all seeders at once using variadic function
	if err := manager.RegisterSeeders([]goseeder.SeederItem{
		{Name: "users", Function: s.UserSeeder},
	}...); err != nil {
		log.Fatalf("Failed to register seeders: %v", err)
	}

	log.Println("All seeders registered successfully")
}
