package main

import (
	"fmt"
	"path/filepath"
)

func potentialEmptyLine() {
	if CLI.Watch == false {
		fmt.Printf("\n")
	}
}

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.LogIfErrFatal(err, "Invalid file path %q", CLI.Path)
	return p
}
