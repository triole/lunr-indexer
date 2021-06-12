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

	chin := make(chan string, CLI.Threads)
	chout := make(chan lunrIndexEntry, CLI.Threads)

	fmt.Printf("\nProcess %d md files, use %d threads\n\n", ln, CLI.Threads)
	bar := progressbar.Default(int64(ln))

	for _, fil := range mdFiles {
		go parseMdFile(fil, chin, chout)
	}

	c := 0
	for li := range chout {
		bar.Add(1)
		lunrIndex = append(lunrIndex, li)
		c++
		if c >= ln {
			close(chin)
			close(chout)
			break
		}
	}

	writeLunrIndexJSON(lunrIndex)
}
