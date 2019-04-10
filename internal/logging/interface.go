package logging

// Facade is common representation of minimal set for logging methods.
type Facade interface {
	// Errorf - should format, append line feed if it is missing and log error message.
	Errorf(format string, a ...interface{})
	// Infof - should format, append line feed if it is missing and log informational message.
	Infof(format string, a ...interface{})
}

// Interface is a logging solution
type Interface interface {
	Facade
	// Close - must close logging interface implementation
	Close() error
}
