package main

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/triole/logseal"
)

func writeLunrIndexJSON(content lunrIndex, outFile string) {
	lg.Info("write lunr index", logseal.F{"path": outFile})
	sort.Sort(lunrIndex(content))
	jsonData, err := json.Marshal(content)
	if err != nil {
		lg.IfErrError(
			"can not marshal lunr index json",
			logseal.F{"content": content, "error": err},
		)
	} else {
		err = os.WriteFile(outFile, jsonData, 0644)
		lg.IfErrError(
			"can not write file",
			logseal.F{"error": err, "path": outFile},
		)
	}
}
