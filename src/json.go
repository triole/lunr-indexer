package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	lg.Log("Write lunr index to %q\n", outFile)
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(outFile, jsonData, 0644)
}
