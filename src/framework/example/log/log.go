package log

import (
	"fmt"
)

type CustomLog struct {
}

// Info logs to the INFO logs.
func (cli CustomLog) Info(args ...interface{}) {
	fmt.Println("Info:", args)
}

func (cli CustomLog) Infof(format string, args ...interface{}) {
	fmt.Println("Infof:", fmt.Sprintf(format, args...))
}

// Warning logs to the WARNING and INFO logs.
func (cli CustomLog) Warning(args ...interface{}) {
	fmt.Println("Warning:", args)
}

func (cli CustomLog) Warningf(format string, args ...interface{}) {
	fmt.Println("Warningf:", fmt.Sprintf(format, args...))
}

// Error logs to the ERROR、 WARNING and INFO logs.
func (cli CustomLog) Error(args ...interface{}) {
	fmt.Println("Error:", args)
}

func (cli CustomLog) Errorf(format string, args ...interface{}) {
	fmt.Println("Error:", fmt.Sprintf(format, args...))
}

// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs.
func (cli CustomLog) Fatal(args ...interface{}) {
	fmt.Println("Fatal:", args)
}
func (cli CustomLog) Fatalf(format string, args ...interface{}) {
	fmt.Println("Fatal:", fmt.Sprintf(format, args...))
}
