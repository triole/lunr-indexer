package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func writeLunrIndexJSON(content []lunrIndexEntry) {
	jsonData, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(CLI.Output, jsonData, 0644)
}
