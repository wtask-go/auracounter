package logging

import (
	"os"
)

type stdout struct {
	*logWriter
}

func (s *stdout) Close() error {
	return nil
}

func NewStdOut(options ...logWriterOption) Interface {
	return &stdout{
		logWriter: buildLogWriter(os.Stdout, options...),
	}
}
