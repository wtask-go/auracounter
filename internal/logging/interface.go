package logging

import "log"

// Facade is common representation of minimal set for logging methods.
type Facade interface {
	// Errorf - must format and log message for error/fatal event.
	Errorfn(string, ...interface{}) error
	// Infof - must format and log message for neutral/notice/warning/information event.
	Infofn(string, ...interface{}) error
}

// Interface is a logging solution
type Interface interface {
	Facade
	// ExposeLogger exposed underlying go-logger. Useful to forward output of external library into loggin.Interface.
	ExposeLogger(prefix string) *log.Logger
	// Close - must close logging interface implementation
	Close() error
}
