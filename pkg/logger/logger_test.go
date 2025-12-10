package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger("debug")
	if log == nil {
		t.Error("Expected logger instance, got nil")
	}
}

func TestLoggerInfo(t *testing.T) {
	log := NewLogger("info")
	// Should not panic
	log.Info("Test info message", "key", "value")
}

func TestLoggerError(t *testing.T) {
	log := NewLogger("error")
	// Should not panic
	log.Error("Test error message", "key", "value")
}

func TestLoggerDebug(t *testing.T) {
	log := NewLogger("debug")
	// Should not panic
	log.Debug("Test debug message", "key", "value")
}

func TestLoggerWarn(t *testing.T) {
	log := NewLogger("warn")
	// Should not panic
	log.Warn("Test warn message", "key", "value")
}
