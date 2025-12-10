package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	AWS      AWSConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Port        string
	LogLevel    string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	DB        int
	KeyPrefix string
}

type JWTConfig struct {
	Secret                string
	RefreshSecret         string
	ExpirationHours       int
	RefreshExpirationDays int
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "SuiteMedia"),
			Environment: getEnv("NODE_ENV", "development"),
			Port:        getEnv("PORT", "3000"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "suitemedia"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 300),
		},
		Redis: RedisConfig{
			Host:      getEnv("REDIS_HOST", "localhost"),
			Port:      getEnvInt("REDIS_PORT", 6379),
			Password:  getEnv("REDIS_PASSWORD", ""),
			DB:        getEnvInt("REDIS_DB", 0),
			KeyPrefix: getEnv("REDIS_KEY_PREFIX", "suitemedia:"),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "your-secret-key"),
			RefreshSecret:         getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			ExpirationHours:       getEnvInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpirationDays: getEnvInt("JWT_REFRESH_EXPIRATION_DAYS", 30),
		},
		CORS: CORSConfig{
			AllowedOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
			AllowedMethods:   strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"), ","),
			AllowedHeaders:   strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-Request-Id"), ","),
			AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		},
	}

	// Validate required fields
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "your-secret-key" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
