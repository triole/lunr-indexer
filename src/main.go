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

	lg = lilog.Init(absPath(CLI.LogFile))

	mdPath := absPath(CLI.Path)
	outJSON := absPath(CLI.Output)

	if CLI.Watch == true {
		watch(mdPath, outJSON)
	} else {
		makeLunrIndex(mdPath, outJSON, CLI.Threads, true)
	}
}

func makeLunrIndex(mdPath string, outFile string, threads int, showProgressBar bool) {
	start := time.Now()

	var bar *progressbar.ProgressBar
	var lunrIndex lunrIndex

	mdFiles := find(mdPath, ".md$")
	ln := len(mdFiles)

	if len(mdFiles) < 1 {
		lg.LogWarn("No md files found in %q\n", mdPath)
	} else {
		chin := make(chan string, threads)
		chout := make(chan lunrIndexEntry, threads)

		potentialEmptyLine()
		lg.Log("No of md files to process %d\n", ln)
		lg.Log("Parallel threads %d\n", threads)
		potentialEmptyLine()

		if showProgressBar == true {
			bar = progressbar.Default(int64(ln))
		}

		for _, fil := range mdFiles {
			go parseMdFile(fil, mdPath, chin, chout)
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
		writeLunrIndexJSON(lunrIndex, outFile)

		lg.Log("Done. It took %s\n", time.Since(start))
		potentialEmptyLine()
	}

}
