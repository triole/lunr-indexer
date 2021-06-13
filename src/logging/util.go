package Logging

import (
	"regexp"
)

func (l Logging) cleanString(s string) string {
	s = l.rxSub(s, "\\x1b[^m]*m", "")
	s = l.rxSub(s, "(\\\\t|\\\\r|\\n)", "")
	return s
}

func (l Logging) rxSub(s string, rx string, repl string) string {
	re := regexp.MustCompile(rx)
	return re.ReplaceAllString(s, repl)
}
