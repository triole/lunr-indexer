package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	BUILDTAGS      string
	appName        = "lunr indexer"
	appDescription = "parse markdown files from a folder and generate a lunr index json"
	appMainversion = "0.1"
)

var CLI struct {
	Path        string `help:"path to scan, default is current dir" arg:"" optional:"" default:"${curdir}"`
	Output      string `help:"json file to write output into" short:"o" default:"${output}"`
	Threads     int    `help:"max threads, default no of avail. cpu threads" short:"t" default:"${proc}"`
	Watch       bool   `help:"watch folder and run rebuild on file change" short:"w"`
	Interval    int32  `help:"watch interval to check for changes in seconds" default:"60" short:"i"`
	Force       bool   `help:"force overwrite of output json file" default:"false" short:"f"`
	LogFile     string `help:"log file" default:"${logfile}" short:"l"`
	VersionFlag bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"curdir":  curdir,
			"logfile": path.Join(os.TempDir(), alnum(appName)+".log"),
			"output":  path.Join(curdir, "lunr-index.json"),
			"proc":    strconv.Itoa(runtime.NumCPU()),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}

func alnum(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile("[^a-z0-9_-]")
	return re.ReplaceAllString(s, "-")
}
