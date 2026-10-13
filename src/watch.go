package main

import (
	"context"
	"regexp"
	"sync"
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

	// Buffered channel prevents blocking when multiple events arrive concurrently
	chin := make(chan time.Time, CLI.Threads)
	itv := time.Duration(CLI.Interval) * time.Second
	go ticker(ctx, chin, itv)
	go runRebuildOnce(ctx, mdPath, outJSON, chin)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chin)
				return
			case <-w.Event:
				chin <- time.Now()
			case err := <-w.Error:
				logger.Fatal("%s", err)
			case <-w.Closed:
				close(chin)
				return
			}
		}
	}()

	if err := w.AddRecursive(mdPath); err != nil {
		logger.Fatal("can not add path to watcher: %v", err)
	}

	if err := w.Start(itv); err != nil {
		logger.Fatal("can not start watcher: %v", err)
	}

	w.Wait()
}

// rebuildState holds the synchronized state for debounced rebuilds
type rebuildState struct {
	mu      sync.Mutex
	timer   *time.Timer
	pending bool
	stopped bool
}

func runRebuildOnce(ctx context.Context, mdPath string, outJSON string, chin chan time.Time) {
	state := &rebuildState{stopped: false}

	for {
		select {
		case <-ctx.Done():
			// Signal shutdown and stop any pending timer
			state.mu.Lock()
			state.stopped = true
			if state.timer != nil {
				state.timer.Stop()
				state.timer = nil
			}
			state.pending = false
			state.mu.Unlock()
			return
		case <-chin:
			// Check if we're already stopping
			state.mu.Lock()
			if state.stopped {
				state.mu.Unlock()
				continue
			}

			// Reset the debounce timer on each event
			if state.pending && state.timer != nil {
				state.timer.Stop()
			}

			state.pending = true
			state.timer = time.AfterFunc(debounceInterval, func() {
				// Lock to check stopped state and update atomically
				state.mu.Lock()
				if state.stopped {
					state.mu.Unlock()
					return
				}

				// Double-check timer hasn't been replaced/stopped
				if state.timer == nil {
					state.mu.Unlock()
					return
				}

				state.pending = false
				state.timer = nil
				state.mu.Unlock()

				// Execute rebuild outside the lock
				makeLunrIndex(mdPath, outJSON, CLI.Threads, false)
			})
			state.mu.Unlock()
		}
	}
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
