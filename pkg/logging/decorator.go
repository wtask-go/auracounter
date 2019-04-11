package logging

import (
	"fmt"
	"unicode"
	"unicode/utf8"
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

// Decorator - formats message before it will be written into log.
// It can includes time/date special formatting, level translation and so on.
//
// `level` - severity level to translate into string if you want;
//
// `message` - source message to write into log;
//
// `idleFrames` - number of runtime frames you want to skip if your decorator adds trace info.
type Decorator func(level SeverityLevel, message string, idleFrames int) string

// severities - default severity levels naming
var severities = [...]string{
	EmergencyLevel: "EMERG",
	AlertLevel:     "ALERT",
	CriticalLevel:  "CRIT",
	ErrorLevel:     "ERR",
	WarningLevel:   "WARNING",
	NoticeLevel:    "NOTICE",
	InfoLevel:      "INFO",
	DebugLevel:     "DEBUG",
}

// defaultDecorator - returns decorator which prepare message like this:
//
// `prefix [YYYY-MM-DD hh:mm:ss.xxxxx] level message`
//
// `prefix` - name of component, app or channel which helps to filter logs in the future.
//
// `timeFormat` - go time format pattern.
//
// `timer` - optional time formatter.
func defaultDecorator(prefix string, timer *Timer) Decorator {
	if prefix != "" && !lastRuneIsSpace(&prefix) {
		prefix += " "
	}
	return func(level SeverityLevel, message string, _ int) string {
		format := "%s[%s] %s"
		if !firstRuneIsSpace(&message) {
			format += " "
		}
		format += "%s"
		return fmt.Sprintf(
			format,
			prefix,
			timer.String(),
			severities[level],
			message,
		)
	}
}

// firstRuneIsSpace - checks the first unicode point in string is space or not.
func firstRuneIsSpace(s *string) bool {
	r, _ := utf8.DecodeRuneInString(*s)
	return unicode.IsSpace(r)
}

// firstRuneIsSpace - checks the last unicode point in string is space or not.
func lastRuneIsSpace(s *string) bool {
	r, _ := utf8.DecodeLastRuneInString(*s)
	return unicode.IsSpace(r)
}
