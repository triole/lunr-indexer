package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	lilog "lunr-indexer/logging"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMakeLunrIndex(t *testing.T) {
	runTest("../testdata/no_md", t)
	runTest("../testdata/set1", t)
}

func runTest(mdFolder string, t *testing.T) {
	logFile := path.Join(os.TempDir(), "lunr-indexer_test.log")
	outFile := path.Join(os.TempDir(), "lunr-indexer_test.json")

	var liAssert lunrIndex
	var li lunrIndex

	lg = lilog.Init(logFile)
	lg.PrintMessages = false
	CLI.Watch = true

	p, _ := filepath.Abs(mdFolder)
	makeLunrIndex(p, outFile, 4, false)

	li = readLunrIndexJSON(outFile)
	assertJSONFile := path.Join(mdFolder, "assert.json")
	if _, err := os.Stat(assertJSONFile); err == nil {
		liAssert = readLunrIndexJSON(assertJSONFile)
	}

	if len(liAssert) > 0 {
		if reflect.DeepEqual(li, liAssert) == false {
			t.Errorf("DeepEqual failed %q != %q", assertJSONFile, outFile)
		}
	}
}

func readLunrIndexJSON(filename string) (li lunrIndex) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("An error occured: %q", err)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(content), &li)
	if err != nil {
		panic(err)
	}
	return li
}
