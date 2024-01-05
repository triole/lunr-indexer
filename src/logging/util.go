package Logging

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

func (l Logging) initContent(str string, itf interface{}) (*logrus.Entry, string) {
	if l.PrintMessages {
		fmt.Printf(str, itf)
	}
	msg := l.cleanString(fmt.Sprintf(str, itf))
	fields := logrus.Fields{"val": itf}
	return l.Logger.WithFields(fields), msg
}

func (l Logging) cleanString(s string) string {
	s = l.rxSub(s, "\\x1b[^m]*m", "")
	s = l.rxSub(s, "(\\\\t|\\\\r|\\n)", "")
	return s
}

func (l Logging) rxSub(s string, rx string, repl string) string {
	re := regexp.MustCompile(rx)
	return re.ReplaceAllString(s, repl)
}
