package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func find(basedir string, rxFilter string) []string {
	if !exists(basedir) {
		logger.Fatal("folder to search inside does not exist: %s", basedir)
	}
	inf, err := os.Stat(basedir)
	if err != nil {
		logger.Fatal("fail to access md folder: %v at %s", err, basedir)
	}
	if !inf.IsDir() {
		logger.Fatal("not a folder, can not look for md files: %s", basedir)
	}

	filelist := []string{}
	rxf, err := regexp.Compile(rxFilter)
	if err != nil {
		logger.Fatal("invalid regex pattern: %v for pattern %s", err, rxFilter)
	}

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			logger.Error("stat file failed: %v at %s", err, path)
			return nil
		}

		// Skip directories
		if f.IsDir() {
			return nil
		}

		// Check for null bytes in path (path traversal defense)
		if strings.Contains(path, "\x00") {
			logger.Warn("skipping file with null byte in path: %s", path)
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

func exists(p string) (b bool) {
	b = true
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		b = false
	}
	return
}
