package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func potentialEmptyLine() {
	if CLI.Watch == false {
		fmt.Printf("\n")
	}
}

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.IfErrFatal(err, "invalid file path %q", CLI.Path)
	return p
}

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	lg.IfErrError(err, "Can not read file %q", filename)
	return
}
