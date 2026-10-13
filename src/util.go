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
		logger.Fatal("invalid path, contains null byte: %s", description)
	}

	// Resolve to absolute path first to prevent traversal attacks
	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		logger.Fatal("invalid path, cannot resolve absolute path: %v", err)
	}

	// Use Lstat to not follow symlinks - reject symlinks for security
	info, err := os.Lstat(resolvedPath)
	if err != nil {
		logger.Fatal("can not access %s %q: %v", description, path, err)
	}

	// Reject symlinks (security: prevents writing/reading via symlink targets)
	if info.Mode()&os.ModeSymlink != 0 {
		logger.Fatal("%s %q is a symbolic link - symlinks not allowed", description, path)
	}

	// Ensure it's a directory (after symlink check, race-condition safe)
	if !info.IsDir() {
		logger.Fatal("%s %q is not a directory", description, path)
	}

	logger.Debug("valid %s: %s, resolved to: %s)", description, path, resolvedPath)
}

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("can not read file %s: %v", filename, err)
		return nil
	}
	return b
}
