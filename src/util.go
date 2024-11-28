package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/triole/logseal"
)

func potentialEmptyLine() {
	if !CLI.Watch {
		fmt.Printf("\n")
	}
}

func absPath(str string) string {
	p, err := filepath.Abs(str)
	lg.IfErrFatal("invalid file path", logseal.F{"error": err, "path": CLI.Path})
	return p
}

func readFile(filename string) (b []byte) {
	b, err := os.ReadFile(filename)
	lg.IfErrError(
		"can not read file %q",
		logseal.F{"error": err, "path": filename},
	)
	return
}
