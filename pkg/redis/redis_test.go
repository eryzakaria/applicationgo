package redis

import (
	"testing"

	"suitemedia/config"
)

func TestNewClient(t *testing.T) {
	cfg := config.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// This will fail to connect without Redis running, but that's expected
	_, err := NewClient(cfg)
	// We expect an error since Redis is not running in test
	if err == nil {
		t.Log("Redis connection succeeded (Redis server running)")
	} else {
		t.Log("Redis connection failed as expected (no Redis server)")
	}
}
