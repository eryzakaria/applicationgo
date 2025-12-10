package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env
	origPort := os.Getenv("SERVER_PORT")
	defer os.Setenv("SERVER_PORT", origPort)

	// Set required environment variables
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_REFRESH_SECRET", "testrefreshsecret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Port is string, should be "8080"
	if cfg.App.Port != "8080" {
		t.Logf("Port check: expected '8080', got '%s'", cfg.App.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected DB host localhost, got %s", cfg.Database.Host)
	}
	if cfg.JWT.Secret != "testsecret" {
		t.Errorf("Expected JWT secret testsecret, got %s", cfg.JWT.Secret)
	}
}

func TestLoadMissingRequiredVar(t *testing.T) {
	// Clear all env vars
	os.Clearenv()

	_, err := Load()
	if err == nil {
		t.Error("Expected error for missing required variables, got nil")
	}
}
