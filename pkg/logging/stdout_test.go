package logging

import (
	"fmt"
	"time"
)

// MakeLog - generates log using Facade implementation, not Interface.
// This func available for all tests.
func MakeLog(f Facade) {
	f.Info("new event", " ", struct{ Foo string }{"Bar"})
	f.Info() // is ignored
	f.Infof("event-1 occurred")
	f.Infof("event #%d occurred", 2)
	f.Error() // is ignored
	f.Error("new event", " ", struct{ Foo string }{"Bar"})
	f.Errorf("error-1 occurred")
	f.Errorf("error #%d occurred", 2)
}

func ExampleNewStdout_withPrefix() {
	logger := NewStdout(
		// blow out current time
		WithDefaultDecoration("test", &Timer{func() time.Time { return time.Time{} }, DefaultTimeFormat}),
	)
	defer logger.Close()
	MakeLog(logger)

	// Output:
	// test [0001-01-01 00:00:00.000000] INFO new event {Bar}
	// test [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// test [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// test [0001-01-01 00:00:00.000000] ERR new event {Bar}
	// test [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// test [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}

func ExampleNewStdout_withoutPrefix() {
	logger := NewStdout(
		// blow out current time
		WithDefaultDecoration("", &Timer{func() time.Time { return time.Time{} }, DefaultTimeFormat}),
	)
	defer logger.Close()
	MakeLog(logger)

	// Output:
	// [0001-01-01 00:00:00.000000] INFO new event {Bar}
	// [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// [0001-01-01 00:00:00.000000] ERR new event {Bar}
	// [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}

func ExampleNewStdout_customTimeFormat() {
	logger := NewStdout(
		// blow out current time
		WithDefaultDecoration(
			"",
			&Timer{
				func() time.Time { return time.Time{}.Add(31 * 24 * time.Hour) }, // 1st february
				"02 Jan 2006 15:04:05",
			},
		),
	)
	defer logger.Close()
	MakeLog(logger)

	// Output:
	// [01 Feb 0001 00:00:00] INFO new event {Bar}
	// [01 Feb 0001 00:00:00] INFO event-1 occurred
	// [01 Feb 0001 00:00:00] INFO event #2 occurred
	// [01 Feb 0001 00:00:00] ERR new event {Bar}
	// [01 Feb 0001 00:00:00] ERR error-1 occurred
	// [01 Feb 0001 00:00:00] ERR error #2 occurred
	//
}

func ExampleNewStdout_customDecoration() {
	logger := NewStdout(
		WithDecoration(
			func(level SeverityLevel, message string, _ int) string {
				severities := [...]string{
					EmergencyLevel: "!!!!",
					AlertLevel:     "!!!",
					CriticalLevel:  "!!",
					ErrorLevel:     "!",
					WarningLevel:   "??",
					NoticeLevel:    "?",
					InfoLevel:      ":",
					DebugLevel:     "@",
				}
				format := "%s [%s] %s"
				return fmt.Sprintf(format, severities[level], "2006-01-02 15:04:05.000000", message)
			},
		),
	)
	defer logger.Close()
	MakeLog(logger)

	// Output:
	// : [2006-01-02 15:04:05.000000] new event {Bar}
	// : [2006-01-02 15:04:05.000000] event-1 occurred
	// : [2006-01-02 15:04:05.000000] event #2 occurred
	// ! [2006-01-02 15:04:05.000000] new event {Bar}
	// ! [2006-01-02 15:04:05.000000] error-1 occurred
	// ! [2006-01-02 15:04:05.000000] error #2 occurred
	//
}
