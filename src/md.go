package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type lunrIndexEntry struct {
	Content  indexEntryContent
	MetaData map[string]interface{}
	URL      string
}

type indexEntryContent struct {
	Md   string
	HTML string
}

func parseMdFile(filename string) (li lunrIndexEntry) {
	var buf bytes.Buffer
	source := readFile(filename)

	markdown := goldmark.New(
		goldmark.WithExtensions(
			goldmarkmeta.Meta,
		),
	)
	context := parser.NewContext()
	if err := markdown.Convert(
		source, &buf, parser.WithContext(context),
	); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)

	li.Content.Md = string(source)
	li.Content.HTML = fmt.Sprintf("%q", buf)
	li.MetaData = metaData

	url := strings.Replace(filename, CLI.Path, "", -1)
	if strings.HasPrefix("/", url) == false {
		url = "/" + url
	}
	li.URL = url
	return
}

func readFile(filename string) (b []byte) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return
}
