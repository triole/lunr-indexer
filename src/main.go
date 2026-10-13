package main

import (
	"fmt"
	"os"
	"sync"
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

	if _, err := os.Stat(outJSON); !os.IsNotExist(err) && !CLI.Force {
		lg.Warn("exit, output json file exists", logseal.F{"path": outJSON})
		fmt.Println("either choose a different output target or use -f/--force to overwrite")
		os.Exit(0)
	}

	if CLI.Watch {
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
	fileCount := len(mdFiles)

	if fileCount < 1 {
		lg.Warn("no md files found", logseal.F{"path": mdPath})
	} else {
		resultChan := make(chan lunrIndexEntry, threads)
		var wg sync.WaitGroup

		potentialEmptyLine()
		lg.Info(
			"process parameters", logseal.F{"no_of_files": fileCount, "threads": threads},
		)
		potentialEmptyLine()

		if showProgressBar {
			bar = progressbar.Default(int64(fileCount))
		}

		// Start worker goroutines for each file
		for _, file := range mdFiles {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				entry := parseMdFileToEntry(f, mdPath)
				resultChan <- entry
			}(file)
		}

		// Close the channel after all goroutines complete
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Collect results
		for entry := range resultChan {
			if showProgressBar {
				bar.Add(1)
			}
			lunrIndex = append(lunrIndex, entry)
		}

		potentialEmptyLine()
		writeLunrIndexJSON(lunrIndex, outFile)

		lg.Info("done", logseal.F{"duration": time.Since(start)})
		potentialEmptyLine()
	}

}
