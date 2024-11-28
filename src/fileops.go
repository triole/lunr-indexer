package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/triole/logseal"
)

func find(basedir string, rxFilter string) []string {
	if !exists(basedir) {
		lg.Fatal(
			"folder to search inside does not exist", logseal.F{"path": basedir},
		)
	}
	inf, err := os.Stat(basedir)
	if err != nil {
		lg.IfErrFatal(
			"fail to access md folder", logseal.F{"error": err, "path": basedir},
		)
	}
	if !inf.IsDir() {
		lg.Fatal(
			"not a folder, can not look for md files", logseal.F{"path": basedir},
		)
	}

	filelist := []string{}
	rxf, _ := regexp.Compile(rxFilter)

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			inf, err := os.Stat(path)
			if err == nil {
				if !inf.IsDir() {
					filelist = append(filelist, path)
				}
			} else {
				lg.IfErrError(
					"stat file failed", logseal.F{"error": err, "path": path},
				)
			}
		}
		return nil
	})
	lg.IfErrFatal(
		"find files failed", logseal.F{"error": err, "path": basedir},
	)
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
