package log

import (
	"configcenter/src/common/blog"
)

// the globle log instance
var logInst Loger = log{}

func init() {
	blog.InitLogs()
}

// SetLoger set a new loger instance
func SetLoger(log Loger) {
	blog.CloseLogs()
	logInst = log
}

// log implements the Loger interface
type log struct{}

// Info logs to the INFO logs.
func (cli log) Info(args ...interface{}) {
	blog.Info("%v", args)
}

// Infof logs to the INFO logs.
func (cli log) Infof(format string, args ...interface{}) {
	blog.Infof(format, args)
}

// Warning logs to the WARNING and INFO logs.
func (cli log) Warning(args ...interface{}) {
	blog.Warn("%v", args)
}

// Warningf logs to the WARNING and INFO logs.
func (cli log) Warningf(format string, args ...interface{}) {
	blog.Warnf(format, args)
}

// Error logs to the Error、 WARNING and INFO logs.
func (cli log) Error(args ...interface{}) {
	blog.Error("%v", args)
}

// Errorf logs to the Error、 WARNING and INFO logs.
func (cli log) Errorf(format string, args ...interface{}) {
	blog.Errorf(format, args)
}

// Fatal logs to the Fatal, Error, WARNING, and INFO logs.
func (cli log) Fatal(args ...interface{}) {
	blog.Fatal(args)
}

// Fatalf logs to the Fatal, Error, WARNING, and INFO logs.
func (cli log) Fatalf(format string, args ...interface{}) {
	blog.Fatalf(format, args)
}

// Info logs to the INFO logs.
func Info(args ...interface{}) {
	logInst.Info(args...)
}

// Infof logs to the INFO logs.
func Infof(format string, args ...interface{}) {
	logInst.Infof(format, args...)
}

// Warning logs to the WARNING and INFO logs.
func Warning(args ...interface{}) {
	logInst.Warning(args...)
}

// Warningf logs to the WARNING and INFO logs.
func Warningf(format string, args ...interface{}) {
	logInst.Warningf(format, args...)
}

// Error logs to the Error、 WARNING and INFO logs.
func Error(args ...interface{}) {
	logInst.Error(args...)
}

// Errorf logs to the Error、 WARNING and INFO logs.
func Errorf(format string, args ...interface{}) {
	logInst.Errorf(format, args...)
}

// Fatal logs to the Fatal, Error, WARNING, and INFO logs.
func Fatal(args ...interface{}) {
	logInst.Fatal(args...)
}

// Fatalf logs to the Fatal, Error, WARNING, and INFO logs.
func Fatalf(format string, args ...interface{}) {
	logInst.Fatalf(format, args...)
}
