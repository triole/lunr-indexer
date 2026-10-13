package main

import (
	"os"
	"path/filepath"
	"regexp"
)

func find(basedir string, rxFilter string) []string {
	validateDir(basedir, "folder to search")

	filelist := []string{}
	rxf, err := regexp.Compile(rxFilter)
	if err != nil {
		logger.Fatal("invalid regex pattern: %v for pattern %s", err, rxFilter)
	}

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			logger.Warn("stat file failed: %v at %s", err, path)
			return nil
		}

		// Skip directories
		if f.IsDir() {
			return nil
		}

		if rxf.MatchString(path) {
			filelist = append(filelist, path)
		}
		return nil
	})
	if err != nil {
		logger.Fatal("find files failed: %v at %s", err, basedir)
	}
	return filelist
}
