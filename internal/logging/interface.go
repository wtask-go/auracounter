package logging

// Facade is common representation of minimal set for logging methods.
type Facade interface {
	// Tracef - must format and log message with stack trace and "trace" tag.
	Tracef(string, ...interface{})
	// Infof - must format and log message for neutral/notice/warning/information event.
	Infof(string, ...interface{})
	// Errorf - must format and log message for error/fatal event.
	Errorf(string, ...interface{})
}

// Interface is a logging solution
type Interface interface {
	Facade

	// Close - must close logging interface implementation
	Close() error
}
