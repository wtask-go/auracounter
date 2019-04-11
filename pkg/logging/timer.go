package logging

import "time"

// Timer - helper type to keep time and its format together.
type Timer struct {
	Now    func() time.Time
	Format string
}

// DefaultTimeFormat - default time format for logging.
const DefaultTimeFormat = "2006-01-02 15:04:05.000000"

// String - Stringer interface implementation for Timer type.
//
// Method works as follow:
//
// * Timer is nil - provides current time in UTC with default time format;
//
// * Timer.Now is nil - provides current time in UTC with Timer.Format;
//
// * provides formatted Timer.Now() time with Timer.Format.
func (t *Timer) String() string {
	switch {
	case t == nil:
		return time.Now().UTC().Format(DefaultTimeFormat)
	case t.Now == nil:
		if t.Format != "" {
			return time.Now().UTC().Format(t.Format)
		}
		return ""
	default:
		if t.Format != "" {
			return t.Now().Format(t.Format)
		}
		return ""
	}
}
