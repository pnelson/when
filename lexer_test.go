package when

import (
	"reflect"
	"testing"
	"time"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		in   string
		want []token
	}{
		// expr
		{
			"",
			[]token{},
		},
		{
			"   ",
			[]token{},
		},
		// lhs (words)
		{
			"a year",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
			},
		},
		{
			"one year",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
			},
		},
		{
			"one year two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year and two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year & two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year, two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year + two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year - two months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"one year two months three weeks four days five hours six minutes seven seconds",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "3"},
				{tokenUnit, "weeks"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "4"},
				{tokenUnit, "days"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "5"},
				{tokenUnit, "hours"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
		{
			"one year two months and three weeks & four days, five hours + six minutes - seven seconds",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "weeks"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "days"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "hours"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
		{
			"one year ago",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenAgo, "ago"},
			},
		},
		{
			"one year before now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenBefore, "before"},
				{tokenNow, "now"},
			},
		},
		{
			"one year after now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenFrom, "after"},
				{tokenNow, "now"},
			},
		},
		{
			"one year from now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenFrom, "from"},
				{tokenNow, "now"},
			},
		},
		// lhs (digit, long unit)
		{
			"1 year",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
			},
		},
		{
			"1 year 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year and 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year & 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year, 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year + 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year - 2 months",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
			},
		},
		{
			"1 year 2 months 3 weeks 4 days 5 hours 6 minutes 7 seconds",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "3"},
				{tokenUnit, "weeks"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "4"},
				{tokenUnit, "days"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "5"},
				{tokenUnit, "hours"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
		{
			"1 year 2 months and 3 weeks & 4 days, 5 hours + 6 minutes - 7 seconds",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "weeks"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "days"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "hours"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
		{
			"1 year 2 months and 3 weeks&4 days,5 hours+6 minutes-7 seconds ",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "months"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "weeks"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "days"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "hours"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
		{
			"1 year ago",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenAgo, "ago"},
			},
		},
		{
			"1 year before now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenBefore, "before"},
				{tokenNow, "now"},
			},
		},
		{
			"1 year after now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenFrom, "after"},
				{tokenNow, "now"},
			},
		},
		{
			"1 year from now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "year"},
				{tokenFrom, "from"},
				{tokenNow, "now"},
			},
		},
		// lhs (digit, short unit)
		{
			"1y",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
			},
		},
		{
			"1y 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y and 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y & 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y, 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y + 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y - 2M",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
			},
		},
		{
			"1y2M3w4d5h6m7s",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "3"},
				{tokenUnit, "w"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "4"},
				{tokenUnit, "d"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "5"},
				{tokenUnit, "h"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "6"},
				{tokenUnit, "m"},
				{tokenOperatorAdd, ""},
				{tokenDigit, "7"},
				{tokenUnit, "s"},
			},
		},
		{
			"1y 2M 3w 4d 5h 6m 7s",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "3"},
				{tokenUnit, "w"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "4"},
				{tokenUnit, "d"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "5"},
				{tokenUnit, "h"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "6"},
				{tokenUnit, "m"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "7"},
				{tokenUnit, "s"},
			},
		},
		{
			"1y 2M and 3w & 4d, 5h + 6m - 7s",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "w"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "d"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "h"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "m"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "s"},
			},
		},
		{
			"1y 2M and 3w&4d,5h+6m-7s ",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "w"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "d"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "h"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "m"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "s"},
			},
		},
		{
			"1y ago",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenAgo, "ago"},
			},
		},
		{
			"1y before now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenBefore, "before"},
				{tokenNow, "now"},
			},
		},
		{
			"1y after now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenFrom, "after"},
				{tokenNow, "now"},
			},
		},
		{
			"1y from now",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenFrom, "from"},
				{tokenNow, "now"},
			},
		},
		// rhs
		{
			"now",
			[]token{
				{tokenNow, "now"},
			},
		},
		{
			"today",
			[]token{
				{tokenDate, "today"},
			},
		},
		{
			"tomorrow",
			[]token{
				{tokenDate, "tomorrow"},
			},
		},
		{
			"yesterday",
			[]token{
				{tokenDate, "yesterday"},
			},
		},
		{
			"midnight",
			[]token{
				{tokenTime, "midnight"},
			},
		},
		{
			"noon",
			[]token{
				{tokenTime, "noon"},
			},
		},
		{
			"3am",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "am"},
			},
		},
		{
			"3pm",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"3 PM",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "PM"},
			},
		},
		{
			"3 in the afternoon",
			[]token{
				{tokenDigit, "3"},
				{tokenKeyword, "in"},
				{tokenKeyword, "the"},
				{tokenKeyword, "afternoon"},
			},
		},
		{
			"3 oclock in the afternoon",
			[]token{
				{tokenDigit, "3"},
				{tokenKeyword, "oclock"},
				{tokenKeyword, "in"},
				{tokenKeyword, "the"},
				{tokenKeyword, "afternoon"},
			},
		},
		{
			"3 o'clock in the afternoon",
			[]token{
				{tokenDigit, "3"},
				{tokenKeyword, "o'clock"},
				{tokenKeyword, "in"},
				{tokenKeyword, "the"},
				{tokenKeyword, "afternoon"},
			},
		},
		{
			"3:04pm",
			[]token{
				{tokenDigit, "3"},
				{tokenColon, ":"},
				{tokenDigit, "04"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"3:04 PM",
			[]token{
				{tokenDigit, "3"},
				{tokenColon, ":"},
				{tokenDigit, "04"},
				{tokenTwelveHour, "PM"},
			},
		},
		{
			"15:04",
			[]token{
				{tokenDigit, "15"},
				{tokenColon, ":"},
				{tokenDigit, "04"},
			},
		},
		{
			"15:04:05",
			[]token{
				{tokenDigit, "15"},
				{tokenColon, ":"},
				{tokenDigit, "04"},
				{tokenColon, ":"},
				{tokenDigit, "05"},
			},
		},
		{
			"quarter to 4pm",
			[]token{
				{tokenKeyword, "quarter"},
				{tokenKeyword, "to"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"quarter after 4pm",
			[]token{
				{tokenKeyword, "quarter"},
				{tokenKeyword, "after"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"quarter past 4pm",
			[]token{
				{tokenKeyword, "quarter"},
				{tokenKeyword, "past"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"half past 4pm",
			[]token{
				{tokenKeyword, "half"},
				{tokenKeyword, "past"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"@3pm",
			[]token{
				{tokenKeyword, "@"},
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"@ 3pm",
			[]token{
				{tokenKeyword, "@"},
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"at 3pm",
			[]token{
				{tokenKeyword, "at"},
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"at noon",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"Sunday",
			[]token{
				{tokenWeekday, time.Sunday.String()},
			},
		},
		{
			"on Wednesday",
			[]token{
				{tokenKeyword, "on"},
				{tokenWeekday, time.Wednesday.String()},
			},
		},
		{
			"January",
			[]token{
				{tokenMonth, time.January.String()},
			},
		},
		{
			"on November",
			[]token{
				{tokenKeyword, "on"},
				{tokenMonth, time.November.String()},
			},
		},
		{
			"1st",
			[]token{
				{tokenDigit, "1"},
				{tokenOrdinal, "st"},
			},
		},
		{
			"2nd",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
			},
		},
		{
			"3rd",
			[]token{
				{tokenDigit, "3"},
				{tokenOrdinal, "rd"},
			},
		},
		{
			"4th",
			[]token{
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
			},
		},
		{
			"4th of the month",
			[]token{
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "of"},
				{tokenKeyword, "the"},
				{tokenUnit, "month"},
			},
		},
		{
			"4th of last month",
			[]token{
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "of"},
				{tokenKeyword, "last"},
				{tokenUnit, "month"},
			},
		},
		{
			"4th of next month",
			[]token{
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "of"},
				{tokenKeyword, "next"},
				{tokenUnit, "month"},
			},
		},
		{
			"last day of March",
			[]token{
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"last day of the month",
			[]token{
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "the"},
				{tokenUnit, "month"},
			},
		},
		{
			"last day of last month",
			[]token{
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "last"},
				{tokenUnit, "month"},
			},
		},
		{
			"last day of next month",
			[]token{
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "next"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last day of March",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"2nd last day of the month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "the"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last day of last month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "last"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last day of next month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenUnit, "day"},
				{tokenKeyword, "of"},
				{tokenKeyword, "next"},
				{tokenUnit, "month"},
			},
		},
		{
			"last Tuesday of March",
			[]token{
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"last Tuesday of the month",
			[]token{
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "the"},
				{tokenUnit, "month"},
			},
		},
		{
			"last Tuesday of last month",
			[]token{
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "last"},
				{tokenUnit, "month"},
			},
		},
		{
			"last Tuesday of next month",
			[]token{
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "next"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last Tuesday of March",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"2nd last Tuesday of the month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "the"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last Tuesday of last month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "last"},
				{tokenUnit, "month"},
			},
		},
		{
			"2nd last Tuesday of next month",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenKeyword, "next"},
				{tokenUnit, "month"},
			},
		},
		{
			"on the 4th",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
			},
		},
		{
			"on the 4th at 4pm",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "at"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"at 4pm on the 4th",
			[]token{
				{tokenKeyword, "at"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
			},
		},
		{
			"on Wednesday at 4pm",
			[]token{
				{tokenKeyword, "on"},
				{tokenWeekday, time.Wednesday.String()},
				{tokenKeyword, "at"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"at 4pm on Wednesday",
			[]token{
				{tokenKeyword, "at"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
				{tokenKeyword, "on"},
				{tokenWeekday, time.Wednesday.String()},
			},
		},
		{
			"on March 14th at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenMonth, time.March.String()},
				{tokenDigit, "14"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"on March the 14th at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "the"},
				{tokenDigit, "14"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"on the 14th March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "14"},
				{tokenOrdinal, "th"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"on the 14th of March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "14"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"at noon on the 14th of March",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "14"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"on the 2nd Tuesday of March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"on the 2nd Tuesday in March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "in"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"at noon on the 2nd Tuesday of March",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"at noon on the 2nd Tuesday in March",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "in"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"on the 2nd last Tuesday of March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"on the 2nd last Tuesday in March at noon",
			[]token{
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "in"},
				{tokenMonth, time.March.String()},
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
			},
		},
		{
			"at noon on the 2nd last Tuesday of March",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
			},
		},
		{
			"at noon on the 2nd last Tuesday in March",
			[]token{
				{tokenKeyword, "at"},
				{tokenTime, "noon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenKeyword, "last"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "in"},
				{tokenMonth, time.March.String()},
			},
		},
		// rhs arithmetic
		{
			"now + 2 days",
			[]token{
				{tokenNow, "now"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "2"},
				{tokenUnit, "days"},
			},
		},
		{
			"now - 2 days",
			[]token{
				{tokenNow, "now"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "2"},
				{tokenUnit, "days"},
			},
		},
		// lhs/rhs
		{
			"7 weeks from Jan 5th at 4pm + 5 days",
			[]token{
				{tokenDigit, "7"},
				{tokenUnit, "weeks"},
				{tokenFrom, "from"},
				{tokenMonth, "Jan"},
				{tokenDigit, "5"},
				{tokenOrdinal, "th"},
				{tokenKeyword, "at"},
				{tokenDigit, "4"},
				{tokenTwelveHour, "pm"},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "5"},
				{tokenUnit, "days"},
			},
		},
		{
			"1y 2M and 3w & 4d, 5h from quarter past 3 o'clock in the afternoon on the 2nd Tuesday of March + 6 minutes - 7 seconds",
			[]token{
				{tokenDigit, "1"},
				{tokenUnit, "y"},
				{tokenOperatorAdd, " "},
				{tokenDigit, "2"},
				{tokenUnit, "M"},
				{tokenOperatorAdd, "and"},
				{tokenDigit, "3"},
				{tokenUnit, "w"},
				{tokenOperatorAdd, "&"},
				{tokenDigit, "4"},
				{tokenUnit, "d"},
				{tokenOperatorAdd, ","},
				{tokenDigit, "5"},
				{tokenUnit, "h"},
				{tokenFrom, "from"},
				{tokenKeyword, "quarter"},
				{tokenKeyword, "past"},
				{tokenDigit, "3"},
				{tokenKeyword, "o'clock"},
				{tokenKeyword, "in"},
				{tokenKeyword, "the"},
				{tokenKeyword, "afternoon"},
				{tokenKeyword, "on"},
				{tokenKeyword, "the"},
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
				{tokenWeekday, time.Tuesday.String()},
				{tokenKeyword, "of"},
				{tokenMonth, time.March.String()},
				{tokenOperatorAdd, "+"},
				{tokenDigit, "6"},
				{tokenUnit, "minutes"},
				{tokenOperatorSub, "-"},
				{tokenDigit, "7"},
				{tokenUnit, "seconds"},
			},
		},
	}
	for _, tt := range tests {
		have, err := lex(tt.in)
		if err != nil {
			t.Fatalf("lex(%q) %v", tt.in, err)
		}
		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("lex(%q)\nhave %v\nwant %v", tt.in, have, tt.want)
		}
	}
}

func TestLexerError(t *testing.T) {
	var tests = []string{
		"/",
		"1year",
		"oneyear",
		"1 ago",
		"one ago",
		"1 from",
		"one from",
		"1 year2 months",
		"one year2 months",
		"1 year2M",
		"one year2M",
	}
	for _, tc := range tests {
		have, err := lex(tc)
		if err == nil {
			t.Errorf("lex(%q)\nhave %v\nwant lex error", tc, have)
		}
	}
}
