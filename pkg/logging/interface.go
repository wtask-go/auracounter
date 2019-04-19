package logging

// Facade is common representation of the set of methods for logging.
type Facade interface {
	// Error - joins arguments with space, append line feed if is missing and log error message.
	Error(a ...interface{})
	// Errorf - format, append line feed if it is missing and log error message.
	Errorf(format string, a ...interface{})

	// Info - joins arguments with space, append line feed if is missing and log informational message.
	Info(a ...interface{})
	// Infof - format, append line feed if it is missing and log informational message.
	Infof(format string, a ...interface{})
}

// Interface contains Facade and is a solution for logging.
type Interface interface {
	Facade
	// Close - must close logging interface implementation
	Close() error
}
