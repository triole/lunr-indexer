package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func potentialEmptyLine() {
	if !CLI.Watch {
		fmt.Printf("\n")
	}
}

func absPath(str string) string {
	p, err := filepath.Abs(str)
	if err != nil {
		logger.Fatal("invalid file path: %v at %s", err, str)
	}
	return p
}

// validatePath checks if a path is safe to use
func validatePath(path string, description string) {
	// Check for null bytes (path traversal attempt)
	if strings.Contains(path, "\x00") {
		logger.Fatal("invalid %s: path contains null byte", description)
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

	logger.Debug("validated %s: %s", description, path)
}

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("can not read file %s: %v", filename, err)
		return nil
	}
	return b
}
