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

type Tags struct {
	Error, Info string
}

// logWriter - wrapper for standart go logger to implement loggin.Interface
type logWriter struct {
	mu             sync.Mutex // protects out
	out            io.Writer  // real writer
	prefix         string
	timeFunc       func() time.Time       // also may define timezone
	timeFormatFunc func(time.Time) string // also may convert time zone
	trace          bool
	tags           Tags
}

type logWriterOption func(*logWriter)

func WithPrefix(prefix string) logWriterOption {
	return func(lw *logWriter) {
		lw.prefix = prefix
	}
}

func WithTimeFunc(f func() time.Time) logWriterOption {
	return func(lw *logWriter) {
		lw.timeFunc = f
	}
}

func WithTimeFormat(f func(time.Time) string) logWriterOption {
	return func(lw *logWriter) {
		lw.timeFormatFunc = f
	}
}

func WithTags(tags Tags) logWriterOption {
	return func(lw *logWriter) {
		lw.tags = tags
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

func buildLogWriter(out io.Writer, options ...logWriterOption) *logWriter {
	lw := &logWriter{
		out:            out,
		timeFunc:       func() time.Time { return time.Now().UTC() },
		timeFormatFunc: func(t time.Time) string { return t.Format("2006-01-02 03:04:05.000000") },
		tags:           Tags{Error: "ERR", Info: "INFO"},
	}
	lw.apply(options...)
	return lw
}

// ExposeLogger - returns go logger which will be write into logging out.
func (lw *logWriter) ExposeLogger(prefix string) *log.Logger {
	if prefix == "" {
		prefix = lw.prefix
	}
	flags := 0
	if lw.trace {
		flags = log.Llongfile
	}
	return log.New(lw.out, prefix, flags)
}

func (lw *logWriter) Write(bytes []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	data := append([]byte(lw.prefix), []byte(lw.timeFormatFunc(lw.timeFunc()))...)
	n, err := lw.out.Write(append(data, bytes...))
	return n, errors.Wrap(err, "logging.stream write error")
}

func (lw *logWriter) getCaller(skip int) (file string, line int) {
	_, file, line, _ = runtime.Caller(skip)
	return file, line
}

func (lw *logWriter) logfn(tag, format string, a ...interface{}) (int, error) {
	if strings.HasSuffix(format, "\n") {
		format = strings.TrimSuffix(format, "\n")
	}
	content := fmt.Sprintf(fmt.Sprintf("%s %s", tag, format), a...)
	if lw.trace {
		f, l := lw.getCaller(3)
		trace := fmt.Sprintf(" (%d:%s)", l, f)
		content += trace
	}
	content += "\n"
	return lw.Write([]byte(content))
}

func (lw *logWriter) Errorfn(format string, a ...interface{}) error {
	_, err := lw.logfn(lw.tags.Error, format, a...)
	return err
}

func (lw *logWriter) Infofn(format string, a ...interface{}) error {
	_, err := lw.logfn(lw.tags.Info, format, a...)
	return err
}
