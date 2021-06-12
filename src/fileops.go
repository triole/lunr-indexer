package main

import (
	"log"
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
	filelist := []string{}
	rxf, _ := regexp.Compile(rxFilter)

	err := filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			filelist = append(filelist, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return filelist
}
