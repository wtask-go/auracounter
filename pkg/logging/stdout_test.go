package logging

import (
	"fmt"
	"time"
)

func ExampleNewStdout_withPrefix() {
	logger := NewStdout(
		WithDecoration(
			// constant zero-time expected
			DefaultVerbosity("test", func() time.Time { return time.Time{} }),
		),
	)
	defer logger.Close()
	logger.Infof("event-1 occurred")
	logger.Infof("event #%d occurred", 2)
	logger.Errorf("error-1 occurred")
	logger.Errorf("error #%d occurred", 2)

	// Output:
	// test [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// test [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// test [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// test [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}

func ExampleNewStdout_withoutPrefix() {
	logger := NewStdout(
		WithDecoration(
			// constant zero-time expected
			DefaultVerbosity("", func() time.Time { return time.Time{} }),
		),
	)
	defer logger.Close()
	logger.Infof("event-1 occurred")
	logger.Infof("event #%d occurred", 2)
	logger.Errorf("error-1 occurred")
	logger.Errorf("error #%d occurred", 2)

	// Output:
	// * [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// * [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// * [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// * [0001-01-01 00:00:00.000000] ERR error #2 occurred
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
	logger.Infof("event-1 occurred")
	logger.Infof("event #%d occurred", 2)
	logger.Errorf("error-1 occurred")
	logger.Errorf("error #%d occurred", 2)

	// Output:
	// : [2006-01-02 15:04:05.000000] event-1 occurred
	// : [2006-01-02 15:04:05.000000] event #2 occurred
	// ! [2006-01-02 15:04:05.000000] error-1 occurred
	// ! [2006-01-02 15:04:05.000000] error #2 occurred
	//
}

func ExampleNewStdout_asFacade() {
	logger := NewStdout(
		WithDecoration(
			// constant zero-time expected
			DefaultVerbosity("test", func() time.Time { return time.Time{} }),
		),
	)
	defer logger.Close()

	func(f Facade) {
		f.Infof("event-1 occurred")
		f.Infof("event #%d occurred", 2)
		f.Errorf("error-1 occurred")
		f.Errorf("error #%d occurred", 2)
	}(logger)

	// Output:
	// test [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// test [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// test [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// test [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}
