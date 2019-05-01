package logging

import (
	"fmt"
	"log"
	"strings"
)

// facade - base unexported type to expose several loggers
type facade struct {
	decorator Decorator
	printer   *log.Logger // is ready for concurrency
}

func (f *facade) println(level SeverityLevel, message string, idleFrames int) {
	if f == nil || f.printer == nil || f.decorator == nil || message == "" {
		// can't print any content without printer and decorator,
		// also method works for nil receiver
		// and why to log empty message?
		return
	}
	if message = f.decorator(level, message, idleFrames); message == "" {
		// the message is completely dropped
		return
	}
	f.printer.Println(message)
}

// Error - joins arguments with space, append line feed if is missing and log error message.
func (f *facade) Error(v ...interface{}) {
	f.println(ErrorLevel, sprint(v...), 3)
}

// Errorf - writes error-level message into log.
func (f *facade) Errorf(format string, v ...interface{}) {
	f.println(ErrorLevel, fmt.Sprintf(format, v...), 3)
}

// Info - joins arguments with space, append line feed if is missing and log informational message.
func (f *facade) Info(v ...interface{}) {
	f.println(InfoLevel, sprint(v...), 3)
}

// Infof - writes informational message into log.
func (f *facade) Infof(format string, v ...interface{}) {
	f.println(InfoLevel, fmt.Sprintf(format, v...), 3)
}

// sprint - formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands. New line IS NOT appended.
func sprint(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return strings.TrimSuffix(fmt.Sprintln(v...), "\n")
}
