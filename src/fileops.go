package main

import (
	"os"
	"path/filepath"
	"regexp"
)

var (
	fi  os.FileInfo
	err error
)

func mkdir(foldername string) {
	os.MkdirAll(foldername, os.ModePerm)
}

func find(basedir string, rxFilter string) []string {
	inf, err := os.Stat(basedir)
	if err != nil {
		lg.LogIfErrFatal(err, "Failed to access md folder %q\n", basedir)
	}
	if inf.IsDir() == false {
		lg.LogFatal(
			"Not a folder %q. Please provide a directory "+
				"to look for md files.\n", basedir,
		)
	}

	filelist := []string{}
	rxf, _ := regexp.Compile(rxFilter)

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			inf, err := os.Stat(path)
			if err == nil {
				if inf.IsDir() == false {
					filelist = append(filelist, path)
				}
			} else {
				lg.LogIfErr(err, "Stat file failed %q", path)
			}
		}
		return nil
	})
	lg.LogIfErrFatal(err, "Find files failed for %q", basedir)
	return filelist
}
