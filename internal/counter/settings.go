package counter

import (
	"errors"
	"fmt"
	"math"
)

type Settings struct {
	StartFrom    int
	Increment    int
	Lower, Upper int
}

func (s *Settings) verify() error {
	if s == nil {
		return errors.New("counter.Settings: unable to verify nil settings")
	}
	if s.Increment < 0 {
		return fmt.Errorf("counter.Settings: negative increment (%d)", s.Increment)
	}
	// Hmm... Zero increment will pause the counter
	// if s.Increment == 0 {
	// 	return errors.New("counter.Settings: useless zero increment")
	// }
	if s.Lower > s.Upper {
		return fmt.Errorf("counter.Settings: invalid counter range [%d:%d]", s.Lower, s.Upper)
	}
	if s.StartFrom < s.Lower || s.StartFrom > s.Upper {
		return fmt.Errorf(
			"counter.Settings: start value (%d) is out of the range [%d:%d]",
			s.StartFrom,
			s.Lower,
			s.Upper,
		)
	}
	if float64(s.Increment) > math.Abs(float64(s.Upper-s.Lower)) {
		return fmt.Errorf(
			"counter.Settings: increment (%d) is wider than counter range [%d:%d]",
			s.Increment,
			s.Lower,
			s.Upper,
		)
	}
	return nil
}

func DefaultSettings() *Settings {
	return &Settings{
		StartFrom: 0,
		Increment: 1,
		Lower:     0,
		Upper:     int(^uint32(0) >> 1), // max int32
	}
}
