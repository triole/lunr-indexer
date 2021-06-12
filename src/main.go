package main

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

func main() {
	parseArgs()

	makeLunrIndex(true)

	if CLI.Watch == true {
		watch()
	}
}

func makeLunrIndex(showProgressBar bool) {
	var bar *progressbar.ProgressBar
	var lunrIndex []lunrIndexEntry

	mdFiles := find(CLI.Path, ".md$")
	ln := len(mdFiles)

	chin := make(chan string, CLI.Threads)
	chout := make(chan lunrIndexEntry, CLI.Threads)

	conditionalPrint(showProgressBar, "\n")
	fmt.Printf("Process %d md files, use %d threads\n", ln, CLI.Threads)
	if showProgressBar == true {
		bar = progressbar.Default(int64(ln))
		fmt.Printf("\n")
	}

	for _, fil := range mdFiles {
		go parseMdFile(fil, chin, chout)
	}

	c := 0
	for li := range chout {
		if showProgressBar == true {
			bar.Add(1)
		}
		lunrIndex = append(lunrIndex, li)
		c++
		if c >= ln {
			close(chin)
			close(chout)
			break
		}
	}

	conditionalPrint(showProgressBar, "\n")
	writeLunrIndexJSON(lunrIndex)
	conditionalPrint(showProgressBar, "\n")
}

func conditionalPrint(b bool, s string) {
	if b == true {
		fmt.Printf(s)
	}
}
