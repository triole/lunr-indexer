package main

import (
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

func watch() {
	w := watcher.New()

	r := regexp.MustCompile(".md$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	chin := make(chan time.Time)
	go ticker(chin)
	go runRebuildOnce(chin)

	go func() {
		for {
			select {
			case _ = <-w.Event:
				chin <- time.Now()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(CLI.Path); err != nil {
		log.Fatalln(err)
	}

	go func() {
		w.Wait()
	}()

	if err := w.Start(time.Duration(CLI.Interval) * time.Second); err != nil {
		log.Fatalln(err)
	}
}

func runRebuildOnce(chin chan time.Time) {
	current := time.Now()
	last := time.Now()
	diff := diffReached(last, current)
	var lastDiff bool
	for t := range chin {
		lastDiff = diff
		last = current
		current = t
		diff = diffReached(last, current)
		if lastDiff == false && diff == true {
			makeLunrIndex(false)
		}
	}
}

func diffReached(last time.Time, current time.Time) bool {
	diff := current.Sub(last)
	return diff > time.Duration(800)*time.Millisecond
}

func ticker(chin chan time.Time) {
	for _ = range time.Tick(time.Duration(1) * time.Second) {
		chin <- time.Now()
	}
}
