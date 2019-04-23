package logging

import (
	"fmt"
	"math/rand"
	"path/filepath"
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

// CustomConstantTimeDecorator - formats row without prefix,
// but with constant timestamp and funny severity level tags.
//
// Expected MakeLog() results :
//
// : [2006-01-02 15:04:05.000000] new event {Bar}
// : [2006-01-02 15:04:05.000000] event-1 occurred
// : [2006-01-02 15:04:05.000000] event #2 occurred
// ! [2006-01-02 15:04:05.000000] new event {Bar}
// ! [2006-01-02 15:04:05.000000] error-1 occurred
// ! [2006-01-02 15:04:05.000000] error #2 occurred
func CustomConstantTimeDecorator() Decorator {
	return func(level SeverityLevel, message string, _ int) string {
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
	}
}

func RandomTestdataName(suffix string) string {
	random := func(l int) string {
		if l <= 0 {
			return ""
		}
		a := []byte("-0123456789_abcdefghijklmnopqrstuvwxyz~")
		s := make([]byte, l)
		for i := range s {
			s[i] = a[rand.Intn(len(a)-1)]
		}
		return string(s)
	}
	return filepath.Join("testdata", random(10)+random(10)+suffix)
}
