package seeders

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/risoftinc/xarch/utils/bcrypt"
)

// User seeding constants
const (
	defaultCreatedBy       = 1
	defaultCreatedEmployee = 100
	defaultAdminRole       = "admin"
	defaultEmployeeRole    = "employee"
)

// User model for seeder
type User struct {
	ID        uint    `gorm:"primaryKey"`
	Username  string  `gorm:"unique;not null"`
	Password  string  `gorm:"not null"`
	Roles     string  `gorm:"not null"`
	Salary    float64 `gorm:"not null"`
	CreatedBy uint
	CreatedAt time.Time
	UpdatedBy *uint
	UpdatedAt *time.Time `gorm:"autoUpdateTime:false"`
}

func (User) TableName() string {
	return "users"
}

// UserSeed seeds user data
func (s *MainSeeder) UserSeed() error {
	log.Println("Seeding users...")

	// Create admin user first
	if err := s.createAdminUser(); err != nil {
		return err
	}

	// Create employee users
	if err := s.createEmployeeUsers(); err != nil {
		return err
	}

	log.Println("Users seeding completed!")
	return nil
}

// createAdminUser creates the admin user
func (s *MainSeeder) createAdminUser() error {
	// Hash password for admin
	hashedPassword, err := bcrypt.HashPasswordWithEnvCost("admin")
	if err != nil {
		return err
	}

	adminUser := User{
		Username:  "admin",
		Password:  string(hashedPassword),
		Roles:     defaultAdminRole,
		Salary:    0,
		CreatedBy: defaultCreatedBy, // Will be updated after admin is created
		CreatedAt: time.Now(),
	}

	// Check if admin already exists
	var count int64
	s.db.Model(&User{}).Where("username = ?", adminUser.Username).Count(&count)

	if count == 0 {
		if err := s.db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Printf("Created admin user: %s (ID: %d)", adminUser.Username, adminUser.ID)

		// Update admin's created_by to reference itself
		if err := s.db.Model(&User{}).Where("id = ?", adminUser.ID).Update("created_by", adminUser.ID).Error; err != nil {
			return err
		}
	} else {
		log.Printf("Admin user already exists: %s", adminUser.Username)
	}

	return nil
}

// createEmployeeUsers creates employee users
func (s *MainSeeder) createEmployeeUsers() error {
	log.Println("Creating", defaultCreatedEmployee, "employee users...")

	// Create new random generator with current time as seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get admin user ID for created_by reference
	var adminUser User
	if err := s.db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using default created_by: %d", defaultCreatedBy)
		adminUser.ID = defaultCreatedBy
	}

	// Salary ranges for employees
	salaryRanges := []struct {
		min, max float64
	}{
		{3000000, 5000000},   // Junior
		{5000000, 8000000},   // Mid-level
		{8000000, 12000000},  // Senior
		{12000000, 15000000}, // Lead/Manager
	}

	// Generate employees
	for i := 1; i <= defaultCreatedEmployee; i++ {
		// Generate fake name
		firstName := faker.FirstName()
		lastName := faker.LastName()

		// Create username from name (lowercase, no spaces)
		username := strings.ToLower(fmt.Sprintf("%s_%s", firstName, lastName))
		username = strings.ReplaceAll(username, " ", "")
		username = strings.ReplaceAll(username, "-", "")
		username = strings.ReplaceAll(username, "'", "")

		// Hash password (same as username)
		hashedPassword, err := bcrypt.HashPasswordWithEnvCost(username)
		if err != nil {
			return err
		}

		// Random salary from ranges
		salaryRange := salaryRanges[i%len(salaryRanges)]
		salary := salaryRange.min + rng.Float64()*(salaryRange.max-salaryRange.min)
		salary = float64(int(salary))

		employeeUser := User{
			Username:  username,
			Password:  string(hashedPassword),
			Roles:     defaultEmployeeRole,
			Salary:    salary,
			CreatedBy: adminUser.ID,
			CreatedAt: time.Now(),
		}

		// Check if user already exists
		var count int64
		s.db.Model(&User{}).Where("username = ?", employeeUser.Username).Count(&count)

		if count == 0 {
			if err := s.db.Create(&employeeUser).Error; err != nil {
				return err
			}
			if i%10 == 0 { // Log every 10th user
				log.Printf("Created employee %d: %s (salary: %.2f)", i, username, salary)
			}
		} else {
			if i%10 == 0 { // Log every 10th user
				log.Printf("Employee already exists: %s", username)
			}
		}
	}

	log.Printf("Created %d employee users", defaultCreatedEmployee)
	return nil
}
