package tai64

import (
	"testing"
	"time"
)

var tai64nTests = []struct {
	hex  string
	time string
}{
	// from `man 8 tai64nlocal`, converted to UTC
	{"@4000000037c219bf2ef02e94", "1999-08-24T04:03:43.7874925Z"},
	// `echo @4000000052c65e550cd675fc | TZ=:/usr/share/zoneinfo/right/Etc/UTC tai64nlocal`
	{"@4000000052c65e550cd675fc", "2014-01-03T06:52:34.2153815Z"},
	// the golang date, converted using http://www.tai64.com/
	{"@4000000043b9410600000000", "2006-01-02T15:04:05Z"},

	// from http://cr.yp.to/libtai/tai64.html
	// the first second of 1970 TAI (which is 10 seconds earlier in UTC)
	{"@400000000000000000000000", "1969-12-31T23:59:50Z"},
	// one second later
	{"@400000000000000100000000", "1969-12-31T23:59:51Z"},
	// 10 seconds later (the first second of 1970 UTC)
	{"@400000000000000A00000000", "1970-01-01T00:00:00Z"},
	// the last second of 1969 TAI (which is 10 seconds earlier in UTC)
	{"@3FFFFFFFFFFFFFFF00000000", "1969-12-31T23:59:49Z"},
	// 1992-06-02 08:07:09 TAI
	{"@400000002a2b2c2d00000000", "1992-06-02T08:06:43Z"},
}

func TestParseTai64n(t *testing.T) {
	for _, test := range tai64nTests {
		result, err := ParseTai64n(test.hex)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if out := result.UTC().Format(time.RFC3339Nano); out != test.time {
			t.Errorf("got %v, expected %v", out, test.time)
		}
	}

	bad := []string{
		// no @
		"4000000037c219bf2ef02e94",
		"4000000037c219bf2ef02e941",
		// too short
		"@4000000037c219bf2ef02e9",
		// too long
		"@4000000037c219bf2ef02e941",
		// too big a number
		"@f000000037c219bf2ef02e94",
		// not hex
		"@G00000000000000000000000",
	}
	for _, test := range bad {
		result, err := ParseTai64n(test)
		if err != ParseError {
			t.Errorf("expected %v, got %v", ParseError, err)
		}
		if !result.IsZero() {
			t.Errorf("expected zero time, got %v", result)
		}
	}
}
