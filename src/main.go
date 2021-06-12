package main

func main() {
	parseArgs()

	mdFiles := find(CLI.Path, ".md$")
	var lunrIndex []lunrIndexEntry

	for _, fil := range mdFiles {
		lunrIndex = append(lunrIndex, parseMdFile(fil))
		break
	}

	writeLunrIndexJSON(lunrIndex)
}
