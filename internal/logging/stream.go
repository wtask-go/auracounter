package logging

import (
	"io"
	"log"
)

type stream struct {
	writer io.WriteCloser
	*facade
}

func (s *stream) Close() error {
	return s.writer.Close()
}

func buildStream(writer io.WriteCloser, options ...facadeOption) Interface {
	return &stream{
		writer: writer,
		facade: (&facade{printer: log.New(writer, "", 0)}).apply(options...),
	}
}
