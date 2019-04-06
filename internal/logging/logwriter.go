package logging

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Tags - enumerates available log-line tags. Pass your choice of tags to WithTags() function.
type Tags struct {
	Error, Info string
}

// logWriter - base unexported type to implement several loggers
type logWriter struct {
	mu             sync.Mutex // protects out
	out            io.Writer  // real writer
	prefix         string
	timeFunc       func() time.Time       // also may define timezone
	timeFormatFunc func(time.Time) string // also may convert time zone
	trace          bool // flag to add or not filename/line num info into log 
	tags           Tags
}

// logWriterOption - func to initialize unexported fields of logWriter struct 
type logWriterOption func(*logWriter)

// WithPrefix - define prefix for every log row.
func WithPrefix(prefix string) logWriterOption {
	return func(lw *logWriter) {
		lw.prefix = prefix
	}
}

// withTimeFunc - define custom time function (for use in tests).
func withTimeFunc(f func() time.Time) logWriterOption {
	return func(lw *logWriter) {
		lw.timeFunc = f
	}
}

// WithTimeFormat - define custom date-time formatter for log row.
func WithTimeFormat(f func(time.Time) string) logWriterOption {
	return func(lw *logWriter) {
		lw.timeFormatFunc = f
	}
}

// WithTags - redefine log tags 
func WithTags(tags Tags) logWriterOption {
	return func(lw *logWriter) {
		lw.tags = tags
	}
}

// WithTrace - add trace info (filename and line number) into every log row.
func WithTrace(trace bool) logWriterOption {
	return func(lw *logWriter) {
		lw.trace = trace
	}
}

func (lw *logWriter) apply(options ...logWriterOption) *logWriter {
	if lw == nil {
		return nil
	}
	for _, o := range options {
		if o != nil {
			o(lw)
		}
	}
	return lw
}

// buildLogWriter - logWriter build helper
func buildLogWriter(out io.Writer, options ...logWriterOption) *logWriter {
	return (&logWriter{
		out:            out,
		timeFunc:       func() time.Time { return time.Now().UTC() },
		timeFormatFunc: func(t time.Time) string { return t.Format("2006-01-02 03:04:05.000000") },
		tags:           Tags{Error: " ERR ", Info: " INFO "},
	}).apply(options...)
}

// ExposeLogger - returns go logger which will be write into logging out.
func (lw *logWriter) ExposeLogger(prefix string) *log.Logger {
	flags := 0
	if lw.trace {
		flags = log.Llongfile
	}
	return log.New(lw, prefix, flags)
}

// Write - implements io.Writer interface to help expose log.Logger.
func (lw *logWriter) Write(bytes []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	log := append([]byte(lw.prefix), []byte(lw.timeFormatFunc(lw.timeFunc()))...)
	n, err := lw.out.Write(append(log, bytes...))
	return n, errors.Wrap(err, "logging.logWriter unable to write")
}

func (lw *logWriter) logf(tag, format string, a ...interface{}) (int, error) {
	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
	}
	content := fmt.Sprintf(fmt.Sprintf("%s%s", tag, format), a...)
	if lw.trace {
		_, file, line, _ := runtime.Caller(2)
		trace := " [%d:%q]"
		if content == "" {
			trace = "[%d:%q]"
		}
		content += fmt.Sprintf(trace, line, file)
	}
	content += "\n"
	return lw.Write([]byte(content))
}

// Errorf - implements logging.Facade interface 
func (lw *logWriter) Errorf(format string, a ...interface{}) error {
	_, err := lw.logf(lw.tags.Error, format, a...)
	return err
}

// Infof - implements logging.Facade interface
func (lw *logWriter) Infof(format string, a ...interface{}) error {
	_, err := lw.logf(lw.tags.Info, format, a...)
	return err
}
