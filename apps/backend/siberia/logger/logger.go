package logger

import (
	"log"
	"os"
)

// Simple logger wrapper for now, can be expanded to use zap/logrus
type Logger struct {
	logger *log.Logger
}

func New(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix+" ", log.LstdFlags),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, v...)
}
