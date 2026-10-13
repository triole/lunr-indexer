package main

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

func watch(mdPath string, outJSON string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := watcher.New()

	r := regexp.MustCompile(".md$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	chin := make(chan time.Time)
	go ticker(chin)
	go runRebuildOnce(mdPath, outJSON, chin)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.Event:
				chin <- time.Now()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				cancel()
				return
			}
		}
	}()

	if err := w.AddRecursive(mdPath); err != nil {
		log.Fatalln(err)
	}

	w.Wait()

	if err := w.Start(time.Duration(CLI.Interval) * time.Second); err != nil {
		log.Fatalln(err)
	}
}

func runRebuildOnce(mdPath string, outJSON string, chin chan time.Time) {
	current := time.Now()
	last := time.Now()
	diff := diffReached(last, current)
	var lastDiff bool
	for t := range chin {
		lastDiff = diff
		last = current
		current = t
		diff = diffReached(last, current)
		if !lastDiff && diff {
			makeLunrIndex(mdPath, outJSON, CLI.Threads, false)
		}
	}
}

func diffReached(last time.Time, current time.Time) bool {
	diff := current.Sub(last)
	return diff > time.Duration(800)*time.Millisecond
}

func ticker(chin chan time.Time) {
	for range time.Tick(time.Duration(1) * time.Second) {
		chin <- time.Now()
	}
}
