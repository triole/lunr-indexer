package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/yuin/goldmark"
)

func parseMdFile(filename string) {
	source := readFile(filename)
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", buf)
}

func readFile(filename string) (b []byte) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return
}
