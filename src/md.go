package main

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/yuin/goldmark"
	goldmarkmeta "github.com/yuin/goldmark-meta"
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
	Title   string `json:"title"`
	Href    string `json:"href"`
	Tags    string `json:"tags"`
	Snippet string `json:"snippet"`
	Content string `json:"content"`
}

// parseMdFileToEntry parses a markdown file and returns a lunrIndexEntry.
func parseMdFileToEntry(filename string, mdPath string) lunrIndexEntry {
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
		logger.Error(
			"parsing markdown failed: %v in %s", err, filename,
		)
		return lunrIndexEntry{}
	}

	metaData := goldmarkmeta.Get(context)
	snippetPlain := makeSnippet(source)
	snippet := markdownToHTML([]byte(snippetPlain))

	href := createHref(filename, mdPath)

	return lunrIndexEntry{
		Title:   getFromMeta("title", metaData, href),
		Href:    href,
		Tags:    makeTags(href),
		Snippet: string(snippet),
		Content: string(source),
	}
}

func getFromMeta(key string, meta map[string]interface{}, alt string) (r string) {
	r = alt
	if val, ok := meta[key]; ok {
		r = val.(string)
	}
	return
}

func makeSnippet(source []byte) (snippet string) {
	const (
		maxSnippetChars = 2000
		maxSnippetLines = 15
	)
	rxHeader, _ := regexp.Compile("^([a-z]+:.*|---)")
	c := 0
	for _, line := range strings.Split(string(source), "\n") {
		nextLine := ""
		if !rxHeader.MatchString(line) {
			nextLine = line + "\n"
		}
		nextChars := len(snippet) + len(nextLine)

		if nextChars > maxSnippetChars || c > maxSnippetLines {
			break
		} else {
			snippet += nextLine
			c++
		}
	}
	return
}

func makeTags(href string) (tags string) {
	// TODO: improve tag generation
	t := strings.Split(href[1:], "/")
	tags = strings.Join(t[0:len(t)-1], ", ")
	return
}

func markdownToHTML(source []byte) (html []byte) {
	html = blackfriday.Run(source)
	return
}

// createHref generates a clean URL path from the filename relative to mdPath.
func createHref(filename string, mdPath string) string {
	href := strings.Replace(filename, mdPath, "", -1)
	href = strings.TrimLeft(href, "/")
	if href != "" {
		href = "/" + href
	}
	return href
}
