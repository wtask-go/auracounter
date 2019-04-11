package logging

import (
	"fmt"
	"log"
)

// facade - base unexported type to expose several loggers
type facade struct {
	decorator Decorator
	printer   *log.Logger // is ready for concurrency
}

func (f *facade) println(level SeverityLevel, message string, idleFrames int) {
	if f == nil || f.printer == nil || f.decorator == nil || message == "" {
		// can't print any content without printer and decorator,
		// so method works for nil receiver
		// why to log empty message?
		return
	}
	if message = f.decorator(level, message, idleFrames); message == "" {
		// the message is completely dropped
		return
	}
	f.printer.Println(message)
}

// Error - joins arguments with space, append line feed if is missing and log error message.
func (f *facade) Error(a ...interface{}) {
	f.println(ErrorLevel, fmt.Sprint(a...), 3)
}

// Errorf - writes error-level message into log.
func (f *facade) Errorf(format string, a ...interface{}) {
	f.println(ErrorLevel, fmt.Sprintf(format, a...), 3)
}

// Info - joins arguments with space, append line feed if is missing and log informational message.
func (f *facade) Info(a ...interface{}) {
	f.println(InfoLevel, fmt.Sprint(a...), 3)
}

// Infof - writes informational message into log.
func (f *facade) Infof(format string, a ...interface{}) {
	f.println(InfoLevel, fmt.Sprintf(format, a...), 3)
}
