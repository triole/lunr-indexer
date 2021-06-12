package main

func main() {
	parseArgs()

	mdFiles := find(CLI.Path, ".md$")
	for _, fil := range mdFiles {
		parseMdFile(fil)
		break
	}
}
