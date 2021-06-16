package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	lg.Log("Write lunr index to %q\n", outFile)
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		lg.LogIfErr(err, "Can not marshal lunr index json %q", content)
	} else {
		err = ioutil.WriteFile(outFile, jsonData, 0644)
		lg.LogIfErr(err, "Can not write file %q\n", outFile)
	}
}
