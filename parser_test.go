package when

import (
	"reflect"
	"testing"
	"time"
)

type testcase struct {
	in   string
	want time.Time
}

func (tc testcase) apply(t *testing.T, now time.Time) {
	t.Helper()
	have, err := ParseNow(tc.in, now)
	if err != nil {
		t.Fatalf("Parse(%q) %v", tc.in, err)
	} else if !reflect.DeepEqual(have, tc.want) {
		t.Errorf("Parse(%q)\nhave %v\nwant %v", tc.in, have, tc.want)
	}
}

func TestParse(t *testing.T) {
	loc, err := time.LoadLocation("MST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Date(2006, time.January, 2, 15, 4, 5, 0, loc)
	tests := []testcase{
		{
			"",
			now,
		},
		{
			"   ",
			now,
		},
		// lhs (words)
		{
			"a year",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year two months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year and two months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year & two months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year, two months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year + two months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year - two months",
			time.Date(2006, time.November, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year two months three weeks four days five hours six minutes seven seconds",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"one year two months and three weeks & four days, five hours + six minutes - seven seconds",
			time.Date(2007, time.March, 27, 20, 9, 58, 0, loc),
		},
		{
			"one year ago",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year before now",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year after now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"one year from now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		// lhs (digit, long unit)
		{
			"1 year",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year 2 months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year and 2 months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year & 2 months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year, 2 months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year + 2 months",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year - 2 months",
			time.Date(2006, time.November, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year 2 months 3 weeks 4 days 5 hours 6 minutes 7 seconds",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"1 year 2 months and 3 weeks & 4 days, 5 hours + 6 minutes - 7 seconds",
			time.Date(2007, time.March, 27, 20, 9, 58, 0, loc),
		},
		{
			"1 year 2 months and 3 weeks&4 days,5 hours+6 minutes-7 seconds ",
			time.Date(2007, time.March, 27, 20, 9, 58, 0, loc),
		},
		{
			"1 year ago",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year before now",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year after now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1 year from now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		// lhs (digit, short unit)
		{
			"1y",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y 2M",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y and 2M",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y & 2M",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y, 2M",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y + 2M",
			time.Date(2007, time.March, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y - 2M",
			time.Date(2006, time.November, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y2M3w4d5h6m7s",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"1y 2M 3w 4d 5h 6m 7s",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"1y 2M and 3w & 4d, 5h + 6m - 7s",
			time.Date(2007, time.March, 27, 20, 9, 58, 0, loc),
		},
		{
			"1y 2M and 3w&4d,5h+6m-7s ",
			time.Date(2007, time.March, 27, 20, 9, 58, 0, loc),
		},
		{
			"1y ago",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y before now",
			time.Date(2005, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y after now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"1y from now",
			time.Date(2007, time.January, 2, 15, 4, 5, 0, loc),
		},
		// rhs
		{
			"now",
			now,
		},
		{
			"today",
			time.Date(2006, time.January, 2, 0, 0, 0, 0, loc),
		},
		{
			"tomorrow",
			time.Date(2006, time.January, 3, 0, 0, 0, 0, loc),
		},
		{
			"yesterday",
			time.Date(2006, time.January, 1, 0, 0, 0, 0, loc),
		},
		{
			"midnight",
			time.Date(2006, time.January, 2, 0, 0, 0, 0, loc),
		},
		{
			"noon",
			time.Date(2006, time.January, 2, 12, 0, 0, 0, loc),
		},
		{
			"2006",
			time.Date(2006, time.January, 1, 0, 0, 0, 0, loc),
		},
		{
			"2006-02",
			time.Date(2006, time.February, 1, 0, 0, 0, 0, loc),
		},
		{
			"2006-01-02",
			time.Date(2006, time.January, 2, 0, 0, 0, 0, loc),
		},
		{
			"2006/02",
			time.Date(2006, time.February, 1, 0, 0, 0, 0, loc),
		},
		{
			"2006/01/02",
			time.Date(2006, time.January, 2, 0, 0, 0, 0, loc),
		},
		{
			"3am",
			time.Date(2006, time.January, 2, 3, 0, 0, 0, loc),
		},
		{
			"3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"3 PM",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"3 in the afternoon",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"3 oclock in the afternoon",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"3 o'clock in the afternoon",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"3:04pm",
			time.Date(2006, time.January, 2, 15, 4, 0, 0, loc),
		},
		{
			"3:04 PM",
			time.Date(2006, time.January, 2, 15, 4, 0, 0, loc),
		},
		{
			"15:04",
			time.Date(2006, time.January, 2, 15, 4, 0, 0, loc),
		},
		{
			"15:04:05",
			time.Date(2006, time.January, 2, 15, 4, 5, 0, loc),
		},
		{
			"quarter to 4pm",
			time.Date(2006, time.January, 2, 15, 45, 0, 0, loc),
		},
		{
			"quarter after 4pm",
			time.Date(2006, time.January, 2, 16, 15, 0, 0, loc),
		},
		{
			"quarter past 4pm",
			time.Date(2006, time.January, 2, 16, 15, 0, 0, loc),
		},
		{
			"half past 4pm",
			time.Date(2006, time.January, 2, 16, 30, 0, 0, loc),
		},
		{
			"@3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"@ 3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at 3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at noon",
			time.Date(2006, time.January, 2, 12, 0, 0, 0, loc),
		},
		{
			"2006-01-02 at 3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"2006/01/02 at 3pm",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at 3pm 2006-01-02",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at 3pm 2006/01/02",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at 3pm on 2006-01-02",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"at 3pm on 2006/01/02",
			time.Date(2006, time.January, 2, 15, 0, 0, 0, loc),
		},
		{
			"Sunday",
			time.Date(2006, time.January, 8, 0, 0, 0, 0, loc),
		},
		{
			"on Wednesday",
			time.Date(2006, time.January, 4, 0, 0, 0, 0, loc),
		},
		{
			"January",
			time.Date(2007, time.January, 1, 0, 0, 0, 0, loc),
		},
		{
			"on November",
			time.Date(2006, time.November, 1, 0, 0, 0, 0, loc),
		},
		{
			"1st",
			time.Date(2006, time.February, 1, 0, 0, 0, 0, loc),
		},
		{
			"2nd",
			time.Date(2006, time.February, 2, 0, 0, 0, 0, loc),
		},
		{
			"3rd",
			time.Date(2006, time.January, 3, 0, 0, 0, 0, loc),
		},
		{
			"4th",
			time.Date(2006, time.January, 4, 0, 0, 0, 0, loc),
		},
		{
			"4th of the month",
			time.Date(2006, time.January, 4, 0, 0, 0, 0, loc),
		},
		{
			"4th of last month",
			time.Date(2005, time.December, 4, 0, 0, 0, 0, loc),
		},
		{
			"4th of next month",
			time.Date(2006, time.February, 4, 0, 0, 0, 0, loc),
		},
		{
			"last day of March",
			time.Date(2006, time.March, 31, 0, 0, 0, 0, loc),
		},
		{
			"last day of the month",
			time.Date(2006, time.January, 31, 0, 0, 0, 0, loc),
		},
		{
			"last day of last month",
			time.Date(2005, time.December, 31, 0, 0, 0, 0, loc),
		},
		{
			"last day of next month",
			time.Date(2006, time.February, 28, 0, 0, 0, 0, loc),
		},
		{
			"2nd last day of March",
			time.Date(2006, time.March, 30, 0, 0, 0, 0, loc),
		},
		{
			"2nd last day of the month",
			time.Date(2006, time.January, 30, 0, 0, 0, 0, loc),
		},
		{
			"2nd last day of last month",
			time.Date(2005, time.December, 30, 0, 0, 0, 0, loc),
		},
		{
			"2nd last day of next month",
			time.Date(2006, time.February, 27, 0, 0, 0, 0, loc),
		},
		{
			"2nd Tuesday of March",
			time.Date(2006, time.March, 14, 0, 0, 0, 0, loc),
		},
		{
			"2nd Tuesday of the month",
			time.Date(2006, time.January, 10, 0, 0, 0, 0, loc),
		},
		{
			"2nd Tuesday of last month",
			time.Date(2005, time.December, 13, 0, 0, 0, 0, loc),
		},
		{
			"2nd Tuesday of next month",
			time.Date(2006, time.February, 14, 0, 0, 0, 0, loc),
		},
		{
			"last Sunday",
			time.Date(2005, time.December, 25, 0, 0, 0, 0, loc),
		},
		{
			"last Monday",
			time.Date(2005, time.December, 26, 0, 0, 0, 0, loc),
		},
		{
			"last Tuesday",
			time.Date(2005, time.December, 27, 0, 0, 0, 0, loc),
		},
		{
			"last Saturday at noon",
			time.Date(2005, time.December, 31, 12, 0, 0, 0, loc),
		},
		{
			"last Tuesday of March",
			time.Date(2006, time.March, 28, 0, 0, 0, 0, loc),
		},
		{
			"last Tuesday of the month",
			time.Date(2006, time.January, 31, 0, 0, 0, 0, loc),
		},
		{
			"last Tuesday of last month",
			time.Date(2005, time.December, 27, 0, 0, 0, 0, loc),
		},
		{
			"last Tuesday of next month",
			time.Date(2006, time.February, 28, 0, 0, 0, 0, loc),
		},
		{
			"2nd last Tuesday of March",
			time.Date(2006, time.March, 21, 0, 0, 0, 0, loc),
		},
		{
			"2nd last Tuesday of the month",
			time.Date(2006, time.January, 24, 0, 0, 0, 0, loc),
		},
		{
			"2nd last Tuesday of last month",
			time.Date(2005, time.December, 20, 0, 0, 0, 0, loc),
		},
		{
			"2nd last Tuesday of next month",
			time.Date(2006, time.February, 21, 0, 0, 0, 0, loc),
		},
		{
			"on the 4th",
			time.Date(2006, time.January, 4, 0, 0, 0, 0, loc),
		},
		{
			"on the 4th at 4pm",
			time.Date(2006, time.January, 4, 16, 0, 0, 0, loc),
		},
		{
			"at 4pm on the 4th",
			time.Date(2006, time.January, 4, 16, 0, 0, 0, loc),
		},
		{
			"on Wednesday at 4pm",
			time.Date(2006, time.January, 4, 16, 0, 0, 0, loc),
		},
		{
			"at 4pm on Wednesday",
			time.Date(2006, time.January, 4, 16, 0, 0, 0, loc),
		},
		{
			"on March 14th at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on March the 14th at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on the 14th March at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on the 14th of March at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"at noon on the 14th of March",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on the 2nd Tuesday of March at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on the 2nd Tuesday in March at noon",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"at noon on the 2nd Tuesday of March",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"at noon on the 2nd Tuesday in March",
			time.Date(2006, time.March, 14, 12, 0, 0, 0, loc),
		},
		{
			"on the 2nd last Tuesday of March at noon",
			time.Date(2006, time.March, 21, 12, 0, 0, 0, loc),
		},
		{
			"on the 2nd last Tuesday in March at noon",
			time.Date(2006, time.March, 21, 12, 0, 0, 0, loc),
		},
		{
			"at noon on the 2nd last Tuesday of March",
			time.Date(2006, time.March, 21, 12, 0, 0, 0, loc),
		},
		{
			"at noon on the 2nd last Tuesday in March",
			time.Date(2006, time.March, 21, 12, 0, 0, 0, loc),
		},
		// rhs arithmetic
		{
			"now + 2 days",
			time.Date(2006, time.January, 4, 15, 4, 5, 0, loc),
		},
		{
			"now - 2 days",
			time.Date(2005, time.December, 31, 15, 4, 5, 0, loc),
		},
		// lhs/rhs
		{
			"7 weeks from Jan 5th at 4pm + 5 days",
			time.Date(2006, time.February, 28, 16, 0, 0, 0, loc),
		},
		{
			"1y 2M and 3w & 4d, 5h from quarter past 3 o'clock in the afternoon on the 2nd Tuesday of March + 6 minutes - 7 seconds",
			time.Date(2007, time.June, 8, 20, 20, 53, 0, loc),
		},
	}
	for _, tc := range tests {
		tc.apply(t, now)
	}
}

func TestParseDST(t *testing.T) {
	loc, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Date(2006, time.January, 2, 15, 4, 5, 0, loc)
	tests := []testcase{
		{
			"6 months",
			time.Date(2006, time.July, 2, 15, 4, 5, 0, loc),
		},
		{
			"now + 6 months",
			time.Date(2006, time.July, 2, 15, 4, 5, 0, loc),
		},
		{
			"6 months before now + 6 months",
			now,
		},
		{
			"1y2M3w4d5h6m7s",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"now + 1y2M3w4d5h6m7s",
			time.Date(2007, time.March, 27, 20, 10, 12, 0, loc),
		},
		{
			"1y2M3w4d5h6m7s before now + 1y2M3w4d5h6m7s",
			now,
		},
	}
	for _, tc := range tests {
		tc.apply(t, now)
	}
}

func TestParserError(t *testing.T) {
	var tests = []string{
		"/",
		"1year",
		"oneyear",
		"1 ago",
		"one ago",
		"1 from",
		"one from",
		"1 year before",
		"one year before",
		"1 year after",
		"one year after",
		"1 year from",
		"one year from",
		"1 year2 months",
		"one year2 months",
		"1 year2M",
		"one year2M",
		"15:04:05 am",
		"3 the in oclock afternoon",
		"on the 3pm",
		"on the 14th noon",
		"on the 14th of the week",
		"on the 14th of last week",
		"on the 14th of next week",
		"on the March the 14th at noon",
		"on the 14th of March the 14th at noon",
		"at noon on March the 14th at 4pm",
		"at noon tomorrow at 4pm",
	}
	now := time.Now()
	for _, tc := range tests {
		have, err := ParseNow(tc, now)
		if err == nil {
			t.Errorf("Parse(%q)\nhave %v\nwant parse error", tc, have)
		}
	}
}
