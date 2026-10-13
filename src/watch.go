package main

import (
	"context"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

const (
	debounceInterval = 800 * time.Millisecond
)

func watch(mdPath string, outJSON string) {
	logger.Info("start watcher, path: %s", mdPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := watcher.New()

	r := regexp.MustCompile(".md$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	chin := make(chan time.Time)
	itv := time.Duration(CLI.Interval) * time.Second
	go ticker(ctx, chin, itv)
	go runRebuildOnce(ctx, mdPath, outJSON, chin)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.Event:
				chin <- time.Now()
			case err := <-w.Error:
				logger.Fatal("%s", err)
			case <-w.Closed:
				cancel()
				return
			}
		}
	}()

	if err := w.AddRecursive(mdPath); err != nil {
		logger.Fatal("%s", err)
	}

	if err := w.Start(itv); err != nil {
		logger.Fatal("%s", err)
	}

	w.Wait()
}

func runRebuildOnce(ctx context.Context, mdPath string, outJSON string, chin chan time.Time) {
	var rebuildTimer *time.Timer
	var rebuildPending bool

	for {
		select {
		case <-ctx.Done():
			if rebuildTimer != nil {
				rebuildTimer.Stop()
			}
			return
		case <-chin:
			// Reset the debounce timer on each event
			if rebuildPending {
				// Cancel pending rebuild and set a new timer
				rebuildTimer.Stop()
			}

			rebuildPending = true
			rebuildTimer = time.AfterFunc(debounceInterval, func() {
				rebuildPending = false
				rebuildTimer = nil
				makeLunrIndex(mdPath, outJSON, CLI.Threads, false)
			})
		}
	}
}

func diffReached(last time.Time, current time.Time) bool {
	diff := current.Sub(last)
	return diff > debounceInterval
}

func ticker(ctx context.Context, chin chan time.Time, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			chin <- time.Now()
		}
	}
}
