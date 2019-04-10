package logging

import (
	"fmt"
	"time"
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
// "*`{prefix}` *`{YYYY-MM-DD hh:mm:ss.xxxxx}` *`{severity level tag}` :`{message}`"
//
// `prefix` - component, app or channel name to help filter logs in the future;
//
// `timer` - optional current time generator, nil-timer is ignored.
// If you need a constant timestamp for log (inside tests, for example)
// or to check time for specific timezone or change date-time format,
// pass a specific timer here. Otherwise, pass nil to use default UTC timer.
func DefaultVerbosity(prefix string, timer func() time.Time) MessageDecorator {
	return func(level SeverityLevel, message string, _ int) string {
		if timer == nil {
			timer = time.Now().UTC
		}
		format := "*%s *%s *%s :%s"
		return fmt.Sprintf(
			format,
			prefix,
			timer().Format("2006-01-02 15:04:05.000000"),
			severities[level],
			message,
		)
	}
}