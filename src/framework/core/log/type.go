package log

// Loger is the interface that must be implemented by every log
type Loger interface {

	// Info logs to the INFO logs.
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	// Warning logs to the WARNING and INFO logs.
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	// Error logs to the ERROR、 WARNING and INFO logs.
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs.
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
