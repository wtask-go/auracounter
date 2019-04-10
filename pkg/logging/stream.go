package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// stream - base type for all streaming loggers
type stream struct {
	*facade
	closer io.Closer
}

// Close - close and free logging stream
func (s *stream) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}

// buildStream - builder for any streming logger.
func buildStream(writer io.Writer, closer io.Closer, options ...facadeOption) Interface {
	return &stream{
		closer: closer,
		facade: (&facade{
			printer:   log.New(writer, "", 0),
			decorator: DefaultVerbosity("", nil),
		}).apply(options...),
	}
}

// NewStdout - creates logger with stdout as writing target.
func NewStdout(options ...facadeOption) Interface {
	return buildStream(os.Stdout, nil, options...)
}

// NewStderr - creates logger with stderr as writing target.
func NewStderr(options ...facadeOption) Interface {
	return buildStream(os.Stderr, nil, options...)
}

// NewNull - creates silent logger without real output.
func NewNull() Interface {
	return buildStream(ioutil.Discard, nil, func(f *facade) { f.decorator = nil })
}

// NewBuffer - creates logger which writes into external buffer. Useful for tests.
func NewBuffer(buffer *[]byte, options ...facadeOption) Interface {
	return nil
}

// NewFile - creates logger which writes log into file.
func NewFile(filename string, options ...facadeOption) Interface {
	return nil
}
