package Logging

import (
	"fmt"
)

// Logf logs an info message
func (l Logging) Logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	str := fmt.Sprintf(msg, a...)
	l.Log.Info(l.cleanString(str))
}

// LogfIfErr logs an error if passed error arg is not nil
func (l Logging) LogfIfErr(err error, msg string, a ...interface{}) {
	if err != nil {
		fmt.Printf(msg, a...)
		str := fmt.Sprintf(msg, a...)
		l.Log.Error(l.cleanString(str))
	}
}
