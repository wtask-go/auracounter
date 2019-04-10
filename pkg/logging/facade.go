package logging

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// SeverityLevel - severity level, used to decorate log rows 
type SeverityLevel int

const (
	// EmergencyLevel - emergency messages.
	// Wikipedia: System is unusable. A panic condition.
	EmergencyLevel SeverityLevel = iota
	// AlertLevel - alert messages.
	// Wikipedia: A condition that should be corrected immediately, such as a corrupted system database.
	AlertLevel
	// CriticalLevel - critical messages. 
	// Wikipedia: Critical conditions. Hard device errors.
	CriticalLevel
	// ErrorLevel - error messages.
	ErrorLevel
	// WarningLevel - warning messages.
	WarningLevel
	// NoticeLevel - notice messages.
	// Wikipedia: Conditions that are not error conditions, but that may require special handling.
	NoticeLevel
	// InfoLevel - informational messages.
	InfoLevel
	// DebugLevel - debug messages.
	// Wikipedia: Messages that contain information normally of use only when debugging a program.
	DebugLevel
)

type (
	// MessageDecorator - finally decorates message before it will be written into log.
	// It can includes time/date special formatting, level translation and so on.
	//
	// `level` - severity level to translate into string if you want; 
	// 
	// `message` - source message to write into log;
	// 
	// `idleFrames` - number of runtime frames you want to skip if your decorator adds trace info.
	MessageDecorator func(level SeverityLevel, message string, idleFrames int) string

	// facade - base unexported type to expose several loggers
	facade struct {
		decorator MessageDecorator
		printer *log.Logger // backend is ready to concurrency
	}

	facadeOption = func(l *facade)
)

func (l *facade) apply(options ...facadeOption) *facade {
	if l == nil {
		return nil
	}
	for _, o := range options {
		if o != nil {
			o(l)
		}
	}
	return l
}

func WithDecoration(d MessageDecorator) facadeOption {
	if d == nil {
		panic(errors.New("logging: can not use nil as MessageDecorator"))
	}
	return func(l *facade) {
		l.decorator = d
	}
}

func (f *facade) println(level SeverityLevel, message string, idleFrames int) {
	if f.decorator == nil {
		// can't print any content without decoration
		return
	}
	message = f.decorator(level, message, idleFrames)
	if message == "" {
		return
	}
	f.printer.Println(message)
}

func (f *facade) Errorf(format string, a ...interface{}) {
	f.println(ErrorLevel, fmt.Sprintf(format, a...), 3)
}

func (f *facade) Infof(format string, a ...interface{}) {
	f.println(InfoLevel, fmt.Sprintf(format, a...), 3)
}
