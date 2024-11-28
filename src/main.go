package main

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()

	lg = logseal.Init(absPath(CLI.LogFile))

	mdPath := absPath(CLI.Path)
	outJSON := absPath(CLI.Output)

	if _, err := os.Stat(outJSON); os.IsNotExist(err) == false && CLI.Force == false {
		lg.Warn("Exitting. Output json file exists %q\n", outJSON)
		fmt.Println("Either choose a different output target or use -f/--force to overwrite")
		os.Exit(0)
	}

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
		lg.Warn("No md files found in %q\n", mdPath)
	} else {
		chin := make(chan string, threads)
		chout := make(chan lunrIndexEntry, threads)

		potentialEmptyLine()
		lg.Info("no of md files to process %d\n", ln)
		lg.Info("parallel threads %d\n", threads)
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

		lg.Info("done, it took %s\n", time.Since(start))
		potentialEmptyLine()
	}

}
