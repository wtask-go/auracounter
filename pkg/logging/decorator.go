package logging

import (
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"
)

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

// DefaultVerbosity - returns decorator which prepare message like this:
//
// `prefix [YYYY-MM-DD hh:mm:ss.xxxxx] level message`
//
// `prefix` - name of component, app or channel which helps to filter logs in the future.
// Empty prefix will be replaced with `*`.
//
// `timer` - optional generator of current time.
// If you need a constant timestamp for log (inside tests, for example)
// or to check time for specific timezone or change date-time format,
// pass a specific timer here. Otherwise, pass nil to use default UTC timer.
func DefaultVerbosity(prefix string, timer func() time.Time) MessageDecorator {
	if prefix == "" {
		prefix = "*"
	}
	if !lastRuneIsSpace(&prefix) {
		prefix += " "
	}
	return func(level SeverityLevel, message string, _ int) string {
		if timer == nil {
			timer = time.Now().UTC
		}
		format := "%s[%s] %s"
		if !firstRuneIsSpace(&message) {
			format += " "
		}
		format += "%s"
		return fmt.Sprintf(
			format,
			prefix,
			timer().Format("2006-01-02 15:04:05.000000"),
			severities[level],
			message,
		)
	}
}

func firstRuneIsSpace(s *string) bool {
	r, _ := utf8.DecodeRuneInString(*s)
	return unicode.IsSpace(r)
}

func lastRuneIsSpace(s *string) bool {
	r, _ := utf8.DecodeLastRuneInString(*s)
	return unicode.IsSpace(r)
}
