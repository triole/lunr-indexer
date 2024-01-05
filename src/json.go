package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	lg.Log("Write lunr index to %q\n", outFile)
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		lg.LogIfErr(err, "Can not marshal lunr index json %q", content)
	} else {
		err = os.WriteFile(outFile, jsonData, 0644)
		lg.LogIfErr(err, "Can not write file %q\n", outFile)
	}
}
