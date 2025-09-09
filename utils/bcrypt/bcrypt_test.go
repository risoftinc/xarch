package bcrypt

import (
	"os"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "special characters password",
			password: "test@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPassword() returned empty hash")
			}
		})
	}
}

func TestHashPasswordWithCost(t *testing.T) {
	tests := []struct {
		name     string
		password string
		cost     int
		wantErr  bool
	}{
		{
			name:     "valid cost",
			password: "testpassword123",
			cost:     10,
			wantErr:  false,
		},
		{
			name:     "minimum cost",
			password: "testpassword123",
			cost:     4,
			wantErr:  false,
		},
		{
			name:     "maximum cost",
			password: "testpassword123",
			cost:     31,
			wantErr:  false,
		},
		{
			name:     "cost below minimum",
			password: "testpassword123",
			cost:     3,
			wantErr:  false,
		},
		{
			name:     "cost above maximum",
			password: "testpassword123",
			cost:     32,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPasswordWithCost(tt.password, tt.cost)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPasswordWithCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPasswordWithCost() returned empty hash")
			}
		})
	}
}

func TestHashPasswordWithEnvCost(t *testing.T) {
	tests := []struct {
		name     string
		password string
		envCost  string
		clearEnv bool
		wantErr  bool
	}{
		{
			name:     "valid env cost",
			password: "testpassword123",
			envCost:  "10",
			wantErr:  false,
		},
		{
			name:     "invalid env cost",
			password: "testpassword123",
			envCost:  "invalid",
			wantErr:  false,
		},
		{
			name:     "empty env cost",
			password: "testpassword123",
			envCost:  "",
			wantErr:  false,
		},
		{
			name:     "no env variable",
			password: "testpassword123",
			clearEnv: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearEnv {
				os.Unsetenv("HASHING_COST")
			} else {
				os.Setenv("HASHING_COST", tt.envCost)
			}

			hash, err := HashPasswordWithEnvCost(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPasswordWithEnvCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("HashPasswordWithEnvCost() returned empty hash")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			password:       "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		minLength int
		want      bool
	}{
		{
			name:      "valid password",
			password:  "testpassword123",
			minLength: 8,
			want:      true,
		},
		{
			name:      "password too short",
			password:  "test",
			minLength: 8,
			want:      false,
		},
		{
			name:      "empty password",
			password:  "",
			minLength: 8,
			want:      false,
		},
		{
			name:      "exact minimum length",
			password:  "12345678",
			minLength: 8,
			want:      true,
		},
		{
			name:      "zero minimum length",
			password:  "test",
			minLength: 0,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.password, tt.minLength); got != tt.want {
				t.Errorf("IsValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     string
	}{
		{
			name:     "weak password",
			password: "12345",
			want:     "weak",
		},
		{
			name:     "medium password",
			password: "123456789",
			want:     "medium",
		},
		{
			name:     "strong password",
			password: "12345678901234567890",
			want:     "strong",
		},
		{
			name:     "empty password",
			password: "",
			want:     "weak",
		},
		{
			name:     "exact weak boundary",
			password: "12345",
			want:     "weak",
		},
		{
			name:     "exact medium boundary",
			password: "123456789",
			want:     "medium",
		},
		{
			name:     "exact strong boundary",
			password: "1234567890",
			want:     "strong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPasswordStrength(tt.password); got != tt.want {
				t.Errorf("GetPasswordStrength() = %v, want %v", got, tt.want)
			}
		})
	}
}
