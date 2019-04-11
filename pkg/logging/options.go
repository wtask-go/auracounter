package logging

import (
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
)

type streamOption = func(s *stream)

// WithDecoration - sets custom decorator which will format log messages.
func WithDecoration(d MessageDecorator) streamOption {
	if d == nil {
		panic(errors.New("logging: can not use nil as MessageDecorator"))
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
// `timer` - optional generator of current time.
// If you need a constant timestamp for log (inside tests, for example)
// or to check time for specific timezone or change date-time format,
// pass a specific timer here. Otherwise, pass nil to use default UTC timer.
func WithDefaultDecoration(prefix string, timer func() time.Time) streamOption {
	return func(s *stream) {
		s.facade.decorator = defaultVerbosity(prefix, timer)
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
