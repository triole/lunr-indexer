package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("can not read file %s: %v", filename, err)
		return nil
	}
	return b
}
