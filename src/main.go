package main

import (
	"os"
	"sync"
	"time"
)

func main() {
	parseArgs()
	defer Close()

	mdPath := absPath(CLI.Path)
	outJSON := absPath(CLI.Output)

	validateDir(mdPath, "input directory")

	if _, err := os.Stat(outJSON); !os.IsNotExist(err) && !CLI.Force {
		logger.Warn("output json file exists: %s, use --force to overwrite", outJSON)
		os.Exit(0)
	}

	if CLI.Watch {
		watch(mdPath, outJSON)
	} else {
		makeLunrIndex(mdPath, outJSON, CLI.Threads)
	}
}

func makeLunrIndex(mdPath string, outFile string, threads int) {
	start := time.Now()

	var lunrIdx lunrIndex

	mdFiles := find(mdPath, ".md$")
	fileCount := len(mdFiles)

	if fileCount < 1 {
		logger.Warn("no md files found in %s", mdPath)
		return
	}

	logger.Info("run indexer: %d files, %d threads", fileCount, threads)

	// Pre-allocate slice for all entries
	entries := make([]lunrIndexEntry, fileCount)

	// Start worker goroutines for each file with panic recovery
	var workersWg sync.WaitGroup
	for i, file := range mdFiles {
		workersWg.Add(1)
		go func(f string, idx int) {
			defer workersWg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic parsing %s: %v", f, r)
					// Leave empty entry on panic
				}
			}()
			entry := parseMdFileToEntry(f, mdPath)
			entries[idx] = entry
		}(file, i)
	}

	// Wait for all workers to finish
	workersWg.Wait()

	// Filter out empty entries
	collectedEntries := make(lunrIndex, 0, fileCount)
	for _, entry := range entries {
		if entry.Href != "" && entry.Title != "" {
			collectedEntries = append(collectedEntries, entry)
		}
	}

	// Use collected entries instead of the parameter-shadowed variable
	lunrIdx = collectedEntries

	writeLunrIndexJSON(lunrIdx, outFile)
	logger.Info("done in %s", time.Since(start))
}
