package logger

import (
	"log"
	"os"
)

type Logger struct {
	level string
}

func NewLogger(level string) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	log.Printf("[FATAL] %s %v", msg, keysAndValues)
	os.Exit(1)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] %s %v", msg, keysAndValues)
	}
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	log.Printf("[WARN] %s %v", msg, keysAndValues)
}
