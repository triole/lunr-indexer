package main

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/triole/logseal"
)

func TestMakeLunrIndex(t *testing.T) {
	runTest("../testdata/no_md", t)
	runTest("../testdata/set1", t)
}

func runTest(mdFolder string, t *testing.T) {
	arr := strings.Split(mdFolder, "/")
	suffix := arr[len(arr)-1]
	logFile := path.Join(os.TempDir(), "lunr-indexer_test_"+suffix+".log")
	outFile := path.Join(os.TempDir(), "lunr-indexer_test_"+suffix+".json")

	var liAssert lunrIndex
	var li lunrIndex

	lg = logseal.Init(logFile)
	CLI.Watch = true

	p, err := filepath.Abs(mdFolder)
	if err != nil {
		t.Errorf("Test failed. Can not set md folder %q", mdFolder)
	}
	if mdFolder != "../testdata/no_md" {
		makeLunrIndex(p, outFile, 4, false)

		li = readLunrIndexJSON(outFile, t)
		assertJSONFile := path.Join(mdFolder, "assert.json")

		if _, err := os.Stat(assertJSONFile); err == nil {
			liAssert = readLunrIndexJSON(assertJSONFile, t)
		}

		if len(liAssert) > 0 {
			if reflect.DeepEqual(li, liAssert) == false {
				t.Errorf("DeepEqual failed %q != %q", assertJSONFile, outFile)
			}
		}
	}
}

func readLunrIndexJSON(filename string, t *testing.T) (li lunrIndex) {
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Can not open file %q", filename)
	} else {
		err = json.Unmarshal([]byte(content), &li)
		if err != nil {
			t.Errorf("Failed to unmarshal %q", filename)
		}
	}
	return li
}
