package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	logger.Info("write lunr index: %s", outFile)
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		logger.Error("can not marshal lunr index json: %v", err)
		return
	}
	err = os.WriteFile(outFile, jsonData, 0644)
	if err != nil {
		logger.Error("can not write file: %v at %s", err, outFile)
	}
}
