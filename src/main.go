package main

import (
	lilog "lunr-indexer/logging"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	lg lilog.Logging
)

func main() {
	parseArgs()
	lg = lilog.Init(CLI.LogFile)

	if CLI.Watch == true {
		watch()
	} else {
		makeLunrIndex(true)
	}
}

func makeLunrIndex(showProgressBar bool) {
	start := time.Now()

	var bar *progressbar.ProgressBar
	var lunrIndex []lunrIndexEntry

	mdFiles := find(CLI.Path, ".md$")
	ln := len(mdFiles)

	chin := make(chan string, CLI.Threads)
	chout := make(chan lunrIndexEntry, CLI.Threads)

	potentialEmptyLine()
	lg.Logf("Process %d md file(s), threads %d\n", ln, CLI.Threads)
	potentialEmptyLine()

	if showProgressBar == true {
		bar = progressbar.Default(int64(ln))
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

	potentialEmptyLine()
	writeLunrIndexJSON(lunrIndex)

	lg.Logf("Done. It took %s\n", time.Since(start))
	potentialEmptyLine()
}
