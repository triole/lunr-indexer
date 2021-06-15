package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type lunrIndex []lunrIndexEntry

func (li lunrIndex) Len() int {
	return len(li)
}

func (li lunrIndex) Less(i, j int) bool {
	return li[i].Href < li[j].Href
}

func (li lunrIndex) Swap(i, j int) {
	li[i], li[j] = li[j], li[i]
}

type lunrIndexEntry struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Tags        string `json:"tags"`
	ContentMD   string `json:"content_md"`
	ContentHTML string `json:"content_html"`
}

type indexEntryMetadata map[string]interface{}

func parseMdFile(filename string, mdPath string, chin chan string, chout chan lunrIndexEntry) {
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
		lg.LogIfErr(err, "Markdown parse fail %q", filename)
		return
	}

	metaData := meta.Get(context)

	href := strings.Replace(filename, mdPath, "", -1)
	re := regexp.MustCompile("^/*")
	href = re.ReplaceAllString(href, "")
	if strings.HasPrefix("/", href) == false {
		href = "/" + href
	}

	// TODO: think about generating better tags
	t := strings.Split(href[1:], "/")
	tags := strings.Join(t[0:len(t)-1], ", ")

	li := lunrIndexEntry{
		// MetaData:    metaData,
		Title:       getFromMeta("title", metaData, href),
		Href:        href,
		Tags:        tags,
		ContentMD:   string(source),
		ContentHTML: fmt.Sprintf("%q", buf),
	}

	chout <- li
	_ = <-chin
}

func readFile(filename string) (b []byte) {
	b, err := ioutil.ReadFile(filename)
	lg.LogIfErr(err, "Can not read file %q", filename)
	return
}

func getFromMeta(key string, meta map[string]interface{}, alt string) (r string) {
	r = alt
	if val, ok := meta[key]; ok {
		r = val.(string)
	}
	return
}
