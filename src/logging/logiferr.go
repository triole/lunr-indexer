package Logging

// LogIfErr log info message if error
func (l Logging) LogIfErr(err error, str string, itf interface{}) {
	if err != nil {
		l.Log(str, itf)
	}
}

// LogIfErrWarn log warn message if error
func (l Logging) LogIfErrWarn(err error, str string, itf interface{}) {
	if err != nil {
		l.LogFatal(str, itf)
	}
}

// LogIfErrFatal log fatal message if err
func (l Logging) LogIfErrFatal(err error, str string, itf interface{}) {
	if err != nil {
		l.LogFatal(str, itf)
	}
}
