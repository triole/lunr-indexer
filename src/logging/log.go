package Logging

import "os"

// Log logs an info message
func (l Logging) Log(str string, itf interface{}) {
	lre, msg := l.initContent(str, itf)
	lre.Info(msg)
}

// LogWarn logs a warning message
func (l Logging) LogWarn(str string, itf interface{}) {
	lre, msg := l.initContent(str, itf)
	lre.Warn(msg)
}

// LogFatal logs a fatal message
func (l Logging) LogFatal(str string, itf interface{}) {
	lre, msg := l.initContent(str, itf)
	lre.Fatal(msg)
	os.Exit(1)
}
