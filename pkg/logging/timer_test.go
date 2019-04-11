package logging

import (
	"testing"
	"time"
)

func TestTimer(test *testing.T) {
	cases := []struct {
		timer         *Timer
		log           string
		checkNotEmpty bool
		parseTime     bool
		expectation   string
		layout        string
		resultType    string
	}{
		{
			// nil Timer is default current time formatter
			(*Timer)(nil), "(*Timer)(nil)", true, true, "", DefaultTimeFormat, "current time",
		},
		{
			&Timer{nil, DefaultTimeFormat}, "&Timer{nil, DefaultTimeFormat}", true, true, "", DefaultTimeFormat, "current time",
		},
		{
			&Timer{nil, ""}, "&Timer{nil, \"\"}", false, false, "", "", "*empty string*",
		},
		{
			// custom func for Now with empty layout
			&Timer{func() time.Time { return time.Date(2019, 4, 11, 22, 18, 0, 0, time.UTC) }, ""},
			"&Timer{time.Date(2019, 4, 11, 22, 18, 0, 0, time.UTC), \"\"}",
			false,
			false,
			"",
			"",
			"*empty string*",
		},
		{
			// custom func for Now with empty layout
			&Timer{func() time.Time { return time.Date(2019, 4, 11, 22, 18, 0, 0, time.UTC) }, DefaultTimeFormat},
			"&Timer{time.Date(2019, 4, 11, 22, 18, 0, 0, time.UTC), DefaultTimeFormat}",
			true,
			true,
			"2019-04-11 22:18:00.000000",
			DefaultTimeFormat,
			"fixed time",
		},
	}

	for i, c := range cases {
		actual := c.timer.String()
		test.Logf("[#%d] %s: %q", i, c.log, actual)
		if c.checkNotEmpty && actual == "" {
			test.Errorf("[#%d] expected %s in %q format, got empty string", i, c.resultType, c.layout)
		}
		if c.parseTime {
			if _, err := time.Parse(c.layout, actual); err != nil {
				test.Errorf("[#%d] failed to parse time (%s): %s", i, c.layout, err)
			}
		}
		if !c.checkNotEmpty && actual != c.expectation {
			test.Errorf("[#%d] expected %q, got %q", i, c.expectation, actual)
		}
	}

}
