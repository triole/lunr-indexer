package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	parseArgs()
	defer Close()

	mdPath := absPath(CLI.Path)
	outJSON := absPath(CLI.Output)

	validatePath(mdPath, "input directory")

	if _, err := os.Stat(outJSON); !os.IsNotExist(err) && !CLI.Force {
		logger.Warn("exit, output json file exists: %s", outJSON)
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

	var lunrIndex lunrIndex

	mdFiles := find(mdPath, ".md$")
	fileCount := len(mdFiles)

	if fileCount < 1 {
		logger.Warn("no md files found in %s", mdPath)
	} else {
		resultChan := make(chan lunrIndexEntry, threads)
		var wg sync.WaitGroup

		logger.Info("run indexer: %d files, %d threads", fileCount, threads)

		// Start worker goroutines for each file with panic recovery
		for _, file := range mdFiles {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				// Recover from any panic to prevent goroutine leak
				defer func() {
					if r := recover(); r != nil {
						logger.Error("panic parsing %s: %v", f, r)
					}
				}()
				entry := parseMdFileToEntry(f, mdPath)
				resultChan <- entry
			}(file)
		}

		// Wait for all workers, then close the channel
		wg.Wait()
		close(resultChan)

		// Collect results (single consumer, no data race)
		for entry := range resultChan {
			// Skip empty entries (from failed parses)
			if entry.Href != "" && entry.Title != "" {
				lunrIndex = append(lunrIndex, entry)
			}
		}
		writeLunrIndexJSON(lunrIndex, outFile)
		logger.Info("done in %s", time.Since(start))
	}

}
