package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
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
	Path        string `help:"path to scan, default is current dir" arg optional default:${curdir}`
	Output      string `help:"json file to write output into" short:o default:${output}`
	VersionFlag bool   `help:"display version" short:V`
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
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag == true {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", " Version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}
