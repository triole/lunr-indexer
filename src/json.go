package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func writeLunrIndexJSON(content []lunrIndexEntry) {
	fmt.Printf("\nWrite lunr index to %q\n", CLI.Output)
	jsonData, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
	}
	_ = ioutil.WriteFile(CLI.Output, jsonData, 0644)
}
