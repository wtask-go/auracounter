package logging

import (
	"bytes"
	"io"
	"os"
)

// stream - base type for all streaming loggers.
type stream struct {
	*facade
	closer io.Closer
}

// Close - close logging stream.
func (s *stream) Close() error {
	if s != nil && s.closer != nil {
		return s.closer.Close()
	}
	return nil
}

// apply - apply given options for stream.
func (s *stream) apply(options ...streamOption) *stream {
	if s == nil {
		return nil
	}
	for _, o := range options {
		if o != nil {
			o(s)
		}
	}
	return s
}

// buildStream - private builder.
func buildStream(options ...streamOption) *stream {
	return (&stream{facade: &facade{}}).apply(options...)
}

// NewStdout - creates logger with stdout as writing target.
//
// Without options logger uses default decoration for log rows:
// `[YYYY-MM-DD hh:mm:ss.xxxxx] severity_tag message`.
func NewStdout(options ...streamOption) Interface {
	// return buildStream(os.Stdout, nil, options...)
	return buildStream(
		withPrintTarget(os.Stdout),
		withCloser(nil),
		WithDefaultDecoration("", nil),
	).apply(options...)
}

// NewStderr - creates logger with stderr as writing target.
//
// Without options logger uses default decoration for log rows:
// `[YYYY-MM-DD hh:mm:ss.xxxxx] severity_tag message`.
func NewStderr(options ...streamOption) Interface {
	return buildStream(
		withPrintTarget(os.Stderr),
		withCloser(nil),
		WithDefaultDecoration("", nil),
	).apply(options...)
}

// NewNull - creates logger without any output.
func NewNull() Interface {
	return (&stream{facade: &facade{}})
}

// NewBuffer - creates logger which writes into external buffer. Useful for tests.
//
// Without options logger uses default decoration for log rows:
// `[YYYY-MM-DD hh:mm:ss.xxxxx] severity_tag message`.
func NewBuffer(buffer *bytes.Buffer, options ...streamOption) Interface {
	return buildStream(
		withPrintTarget(buffer),
		withCloser(nil),
		WithDefaultDecoration("", nil),
	).apply(options...)
}

// NewFile - creates logger which writes log into file.
func NewFile(filename string, options ...streamOption) Interface {
	return nil // TODO implement print to file
}
