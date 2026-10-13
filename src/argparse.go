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
	Interval    int64  `help:"watch interval to check for changes in seconds" default:"60" short:"i"`
	Force       bool   `help:"force overwrite of output json file" default:"false" short:"f"`
	LogLevel    string `help:"log level: debug, info, warn, or error" short:"v" default:"info"`
	LogFile     string `help:"log file path" short:"l"`
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
			"curdir": curdir,
			"output": path.Join(curdir, "lunr-index.json"),
			"proc":   strconv.Itoa(runtime.NumCPU()),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}

	SetLogLevel(ParseLogLevel(CLI.LogLevel))

	if err := SetLogFile(CLI.LogFile); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] failed to open log file: %v\n", err)
		os.Exit(1)
	}
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.ReplaceAll(s, "_subversion: ", "Version: "+appMainversion+".")
	fmt.Printf("%s\n", s)
}

func alnum(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile("[^a-z0-9_-]")
	return re.ReplaceAllString(s, "-")
}
