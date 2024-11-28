package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	lg.Info("write lunr index to %q\n", outFile)
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		lg.IfErrError(err, "can not marshal lunr index json %q", content)
	} else {
		err = os.WriteFile(outFile, jsonData, 0644)
		lg.IfErrError(err, "can not write file %q\n", outFile)
	}
}
