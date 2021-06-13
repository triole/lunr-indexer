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
	MetaData indexEntryMetadata
	URL      string
}

type indexEntryContent struct {
	Md   string
	HTML string
}

type indexEntryMetadata map[string]interface{}

func parseMdFile(filename string, chin chan string, chout chan lunrIndexEntry) {
	chin <- filename

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
		lg.LogfIfErr(err, "Markdown parse fail %q", filename)
		return
	}

	metaData := meta.Get(context)

	url := strings.Replace(filename, CLI.Path, "", -1)
	if strings.HasPrefix("/", url) == false {
		url = "/" + url
	}

	li := lunrIndexEntry{
		Content: indexEntryContent{
			Md:   string(source),
			HTML: fmt.Sprintf("%q", buf),
		},
		MetaData: metaData,
		URL:      url,
	}

	chout <- li
	_ = <-chin
}

func readFile(filename string) (b []byte) {
	b, err := ioutil.ReadFile(filename)
	lg.LogfIfErr(err, "Can not read file %q", filename)
	return
}
