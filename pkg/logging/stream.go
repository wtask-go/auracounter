package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type stream struct {
	*facade
	closer io.Closer
}

func (s *stream) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}

func buildStream(writer io.Writer, closer io.Closer, options ...facadeOption) Interface {
	return &stream{
		closer: closer,
		facade: (&facade{
			printer:   log.New(writer, "", 0),
			decorator: DefaultVerbosity("", nil),
		}).apply(options...),
	}
}

func NewStdout(options ...facadeOption) Interface {
	return buildStream(os.Stdout, nil, options...)
}

func NewStderr(options ...facadeOption) Interface {
	return buildStream(os.Stderr, nil, options...)
}

func NewNull() Interface {
	return buildStream(ioutil.Discard, nil, func(f *facade) { f.decorator = nil })
}

func NewBuffer(buffer *[]byte, options ...facadeOption) Interface {
	return nil
}

func NewFile(filename string, options ...facadeOption) Interface {
	return nil
}
