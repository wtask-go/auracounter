package logging

import (
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"
)

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

// defaultVerbosity - returns decorator which prepare message like this:
//
// `prefix [YYYY-MM-DD hh:mm:ss.xxxxx] level message`
//
// `prefix` - name of component, app or channel which helps to filter logs in the future.
//
// `timer` - optional generator of current time.
// If you need a constant timestamp for log (inside tests, for example)
// or to check time for specific timezone or change date-time format,
// pass a specific timer here. Otherwise, pass nil to use default UTC timer.
func defaultVerbosity(prefix string, timer func() time.Time) MessageDecorator {
	if prefix != "" && !lastRuneIsSpace(&prefix) {
		prefix += " "
	}
	if timer == nil {
		timer = time.Now().UTC
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
			timer().Format("2006-01-02 15:04:05.000000"),
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
