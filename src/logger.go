package main

import (
	"log"
	"os"
	"strings"
	"sync"
)

// Logger handles application logging with levels
type Logger struct {
	logger *log.Logger
	level  LogLevel
	mu     sync.Mutex
	file   *os.File
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

var logger = &Logger{
	logger: log.New(os.Stderr, "", 0),
	level:  LevelInfo,
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[INFO]  "+format+"\n", args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[WARN]  "+format+"\n", args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[ERROR] "+format+"\n", args...)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Fatalf("[FATAL] "+format+"\n", args...)
}

func SetLogFile(path string) error {

	if path == "" {
		return nil
	}

	// Close existing file if open
	if logger.file != nil {
		logger.file.Close()
	}

	// Open with append mode - creates file if it doesn't exist
	// O_CREATE: create if not exists
	// O_WRONLY: write-only
	// O_APPEND: append to end of file instead of overwriting
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.file = file

	// Create a new logger with file as output
	logger.logger = log.New(file, "", log.LstdFlags)

	return nil
}

func SetLogLevel(level LogLevel) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.level = level
}

// ParseLogLevel converts a string to LogLevel
func ParseLogLevel(s string) LogLevel {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo // default fallback
	}
}

// Close closes the log file if it's open
func Close() {
	if logger.file != nil {
		logger.file.Close()
		logger.file = nil
	}
}
