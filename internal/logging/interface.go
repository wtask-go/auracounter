package logging

import "log"

// Facade is common representation of minimal set for logging methods.
type Facade interface {
	// Errorf - must format, tag as "error" and log message.
	Errorf(string, ...interface{}) error
	// Infof - must format, tag as "information" and log message.
	Infof(string, ...interface{}) error
}

// Interface is a logging solution
type Interface interface {
	Facade
	// ExposeLogger - builds go-logger which will forward all logs into loggin.Interface.
	ExposeLogger(prefix string) *log.Logger
	// Close - must close logging interface implementation
	Close() error
}
