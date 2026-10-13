package main

import (
	"os"
	"path/filepath"
	"strings"
)

func absPath(str string) string {
	p, err := filepath.Abs(str)
	if err != nil {
		logger.Fatal("invalid file path: %v at %s", err, str)
	}
	return p
}

// validateDir checks if a path is safe to use
// Calls logger.Fatal and exits if validation fails
func validateDir(path string, description string) {
	// Check for null bytes (path traversal attempt)
	if strings.Contains(path, "\x00") {
		logger.Fatal("invalid path, contains null byte: %s:", description)
	}

	// Check path exists and is accessible
	info, err := os.Stat(path)
	if err != nil {
		logger.Fatal("can not access %s '%s': %v", description, path, err)
	}

	// Ensure it's a directory
	if !info.IsDir() {
		logger.Fatal("%s '%s' is not a directory", description, path)
	}

	logger.Debug("valid %s: %s", description, path)
}

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("can not read file %s: %v", filename, err)
		return nil
	}
	return b
}
