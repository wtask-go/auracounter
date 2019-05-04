package counter

import (
	"strings"
	"testing"
)

func TestSettingsVerification(t *testing.T) {
	cases := []struct {
		s      *Settings
		errMsg string
	}{
		{nil, "unable to verify nil settings"},
		{&Settings{}, ""},
		{&Settings{Increment: -1}, "negative increment (-1)"},
		{&Settings{Increment: 1}, "increment (1) is wider than counter range [0:0]"},
		{&Settings{Increment: 10}, "increment (10) is wider than counter range [0:0]"},
		{&Settings{Increment: 10, Upper: 9}, "increment (10) is wider than counter range [0:9]"},
		{&Settings{Increment: 10, Upper: 10}, ""},
		{&Settings{Increment: 0, Lower: 1}, "invalid counter range [1:0]"},
		{&Settings{Increment: 0, Lower: -1, Upper: -10}, "invalid counter range [-1:-10]"},
		{&Settings{Increment: 0, Lower: -1, Upper: 1}, ""},
		{&Settings{Increment: 1, Lower: -1, Upper: 1}, ""},
		{&Settings{Increment: 2, Lower: -1, Upper: 1}, ""},
		{&Settings{StartFrom: 10}, "start value (10) is out of the range [0:0]"},
		{&Settings{StartFrom: -10}, "start value (-10) is out of the range [0:0]"},
		{&Settings{StartFrom: 1, Upper: 10}, ""},
	}

	for _, c := range cases {
		err := c.s.verify()

		switch {
		case err != nil && c.errMsg == "":
			t.Errorf("Unexpected error: %v", err)
		case err != nil && c.errMsg != "":
			if !strings.Contains(err.Error(), c.errMsg) {
				t.Errorf("Expected error will contain %q, but got: %v", c.errMsg, err)
			}
		case err == nil && c.errMsg != "":
			t.Errorf("Expected error will contain %q, but got nothing", c.errMsg)
		}
	}
}
