package logging

import (
	"io"
	"log"

	"github.com/pkg/errors"
)

type streamOption = func(s *stream)

// WithDecoration - sets custom decorator which will format log messages.
func WithDecoration(d Decorator) streamOption {
	if d == nil {
		panic(errors.New("can not use nil as logging.Decorator"))
	}
	return func(s *stream) {
		s.facade.decorator = d
	}
}

// WithDefaultDecoration - provide default formatting for log rows like this:
//
// `prefix [YYYY-MM-DD hh:mm:ss.xxxxx] level message`
//
// `prefix` - name of component, app or channel which helps to filter logs in the future.
//
// `timer` - optional time formatter.
// Check Timer.String() method docs to know how Timer formats time.
func WithDefaultDecoration(prefix string, timer *Timer) streamOption {
	return func(s *stream) {
		s.facade.decorator = defaultDecorator(prefix, timer)
	}
}

// withPrintTarget - private option to init facade.printer
func withPrintTarget(writer io.Writer) streamOption {
	return func(s *stream) {
		s.facade.printer = log.New(writer, "", 0)
	}
}

// withCloser - private option to init stream.closer
func withCloser(closer io.Closer) streamOption {
	return func(s *stream) {
		s.closer = closer
	}
}
