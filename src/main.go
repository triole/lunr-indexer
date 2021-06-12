package main

import "fmt"

func main() {
	parseArgs()

	mdFiles := find(CLI.Path, ".md$")
	var lunrIndex []lunrIndexEntry

	for _, fil := range mdFiles {
		lunrIndex = append(lunrIndex, parseMdFile(fil))
		break
	}

	fmt.Printf("%+v\n", lunrIndex)
}
