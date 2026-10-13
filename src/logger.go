package main

import (
	"log"
	"os"
	"sync"
)

// Logger handles application logging with levels
type Logger struct {
	logger *log.Logger
	level  LogLevel
	mu     sync.Mutex
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

func SetLogLevel(level LogLevel) {
	logger.level = level
}
