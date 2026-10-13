package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger handles application logging with levels
type Logger struct {
	level   LogLevel
	mu      sync.Mutex
	file    *os.File
	useFile bool // tracks if we're using file output vs direct stdout/stderr
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

var logger = &Logger{
	level:   LevelInfo,
	useFile: false,
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.write("[DEBUG]", format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.write("[INFO]", format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.write("[WARN]", format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.write("[ERROR]", format, args...)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.write("[FATAL]", format, args...)
	os.Exit(1)
}

// write handles the actual output to file or stdout/stderr
func (l *Logger) write(prefix string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000 MST")
	message := fmt.Sprintf(format, args...)
	output := fmt.Sprintf("%s %s %s\n", timestamp, prefix, message)
	if l.useFile && l.file != nil {
		fmt.Fprint(l.file, output)
	} else {
		fmt.Fprint(os.Stderr, output)
	}
}

func SetLogFile(path string) error {
	if path == "" {
		// Close file if previously opened
		if logger.file != nil {
			logger.file.Close()
			logger.file = nil
			logger.useFile = false
		}
		return nil
	}

	// Close existing file if open
	if logger.file != nil {
		logger.file.Close()
	}

	// Open with append mode - creates file if it doesn't exist
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.file = file
	logger.useFile = true
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
		logger.useFile = false
	}
}
