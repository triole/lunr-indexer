package Logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logging contains all log output operations
type Logging struct {
	Logger        *logrus.Logger
	LogFile       string
	LogFileOS     *os.File
	PrintMessages bool
}

// Init does exactly what it says, initializing the Logging class
func Init(logFile string) (l Logging) {
	var err error
	var log = logrus.New()
	l.LogFile = logFile
	l.PrintMessages = true

	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "date",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
		},
		TimestampFormat: "2006-01-02 15:04:05.000 MST",
		PrettyPrint:     false,
	})

	l.LogFileOS, err = os.OpenFile(
		l.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666,
	)
	if err == nil {
		log.Out = l.LogFileOS
	} else {
		log.Info("Failed to log to file, use default stderr")
	}

	l.Logger = log
	return
}
