package bcrypt

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Generate hash with default cost (10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// HashPasswordWithCost hashes a password using bcrypt with custom cost
func HashPasswordWithCost(password string, cost int) (string, error) {
	// Validate cost range
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// HashPasswordWithEnvCost hashes a password using bcrypt with cost from HASHING_COST env var
func HashPasswordWithEnvCost(password string) (string, error) {
	// Get cost from environment variable
	costStr := os.Getenv("HASHING_COST")

	// If HASHING_COST is not set or empty, use default cost
	if costStr == "" {
		return HashPassword(password)
	}

	// Parse cost from string to int
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		// If parsing fails, use default cost
		return HashPassword(password)
	}

	// Use HashPasswordWithCost with parsed cost
	return HashPasswordWithCost(password, cost)
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsValidPassword checks if a password is valid (not empty and has minimum length)
func IsValidPassword(password string, minLength int) bool {
	if password == "" || len(password) < minLength {
		return false
	}
	return true
}

// GetPasswordStrength returns password strength level
func GetPasswordStrength(password string) string {
	if len(password) < 6 {
		return "weak"
	} else if len(password) < 10 {
		return "medium"
	} else {
		return "strong"
	}
}
