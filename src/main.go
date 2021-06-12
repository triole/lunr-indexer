package main

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

func main() {
	var lunrIndex []lunrIndexEntry

	parseArgs()

	mdFiles := find(CLI.Path, ".md$")
	ln := len(mdFiles)

	fmt.Printf("Process %d md files\n", ln)
	bar := progressbar.Default(int64(ln))
	for _, fil := range mdFiles {
		bar.Add(1)
		lunrIndex = append(lunrIndex, parseMdFile(fil))
	}

	writeLunrIndexJSON(lunrIndex)
}
