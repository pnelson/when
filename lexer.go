package when

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenAgo tokenType = iota
	tokenBefore
	tokenColon
	tokenDate
	tokenDigit
	tokenEOF
	tokenError
	tokenFrom
	tokenKeyword
	tokenMonth
	tokenNow
	tokenOperatorAdd
	tokenOperatorSub
	tokenOrdinal
	tokenTime
	tokenTwelveHour
	tokenUnit
	tokenWeekday
)

const eof = rune(-1)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenError:
		return t.val
	case tokenEOF:
		return "EOF"
	}
	return fmt.Sprintf("%q", t.val)
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input  string
	i, j   int // position within input
	width  int // width of last rune
	tokens []token
}

func lex(s string) ([]token, error) {
	l := &lexer{
		input:  s,
		tokens: make([]token, 0),
	}
	for state := readExpr; state != nil; {
		state = state(l)
	}
	if len(l.tokens) > 0 {
		last := l.tokens[len(l.tokens)-1]
		if last.typ == tokenError {
			return nil, errors.New(last.val)
		}
	}
	return l.tokens, nil
}

func (l *lexer) emit(typ tokenType) {
	v := l.value()
	l.emitAs(typ, v)
}

func (l *lexer) emitAs(typ tokenType, value string) {
	l.tokens = append(l.tokens, token{typ, value})
	l.i = l.j
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens = append(l.tokens, token{tokenError, fmt.Sprintf(format, args...)})
	return nil
}

func (l *lexer) ignore() {
	l.i = l.j
}

func (l *lexer) read() rune {
	if l.j >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.j:])
	l.j += width
	l.width = width
	return r
}

func (l *lexer) readFn(fn func(rune) bool) {
	for {
		r := l.read()
		if r == eof || !fn(r) {
			l.unread()
			break
		}
	}
}

func (l *lexer) readRune(ch rune) bool {
	r := l.read()
	return r == ch
}

func (l *lexer) readString(s string) bool {
	for _, r := range s {
		if !l.readRune(r) {
			l.unread()
			return false
		}
	}
	r := l.peek()
	return unicode.IsSpace(r) || r == eof
}

func (l *lexer) peek() rune {
	r := l.read()
	l.unread()
	return r
}

func (l *lexer) peekFn(fn func(rune) bool) string {
	s := ""
	width := l.width
	for {
		r := l.read()
		if r == eof || !fn(r) {
			l.unread()
			break
		}
		s += string(r)
	}
	for i := range s {
		_, width := utf8.DecodeRuneInString(s[i:])
		l.j -= width
	}
	l.width = width
	return s
}

func (l *lexer) unread() {
	l.j -= l.width
}

func (l *lexer) value() string {
	return l.input[l.i:l.j]
}

func readExpr(l *lexer) stateFn {
	l.readFn(unicode.IsSpace)
	l.ignore()
	r := l.peek()
	switch {
	case r == eof:
		return nil
	case r == '@':
		return readAtSymbol
	case r == '+':
		l.read()
		l.emit(tokenOperatorAdd)
		return readExpr
	case r == '-':
		l.read()
		l.emit(tokenOperatorSub)
		return readExpr
	case unicode.IsDigit(r):
		return readDigit
	case unicode.IsLetter(r):
		return readLetter
	}
	return l.errorf("invalid character")
}

func readAtSymbol(l *lexer) stateFn {
	ok := l.readRune('@')
	if !ok {
		return l.errorf("invalid character")
	}
	l.emit(tokenKeyword)
	return readExpr
}

func readColon(l *lexer) stateFn {
	l.read()
	l.emit(tokenColon)
	r := l.peek()
	if !unicode.IsDigit(r) {
		return l.errorf("colon must be followed by a digit")
	}
	return readDigit
}

func readDigit(l *lexer) stateFn {
	l.readFn(unicode.IsDigit)
	l.emit(tokenDigit)
	r := l.peek()
	switch r {
	case eof:
		return nil
	case 's':
		return readDurationUnitSecondsOrOrdinal
	case 'm', 'h', 'd', 'w', 'M', 'y':
		l.read()
		l.emit(tokenUnit)
		return readDurationNextShort
	case 'n', 'r', 't':
		return readOrdinal
	case 'a', 'p':
		return readTwelveHour
	case ':':
		return readColon
	}
	return readExpr
}

func readDurationNext(l *lexer) stateFn {
	r := l.peek()
	switch r {
	case eof:
		return nil
	case ',':
		fallthrough
	case '+':
		fallthrough
	case '&':
		l.read()
		l.emit(tokenOperatorAdd)
		return readExpr
	case '-':
		l.read()
		l.emit(tokenOperatorSub)
		return readExpr
	}
	if !unicode.IsSpace(r) {
		return l.errorf("invalid character")
	}
	return readDurationSpace
}

func readDurationNextShort(l *lexer) stateFn {
	r := l.peek()
	switch {
	case r == eof:
		return nil
	case r == ',':
		fallthrough
	case r == '+':
		fallthrough
	case r == '&':
		l.read()
		l.emit(tokenOperatorAdd)
		return readExpr
	case r == '-':
		l.read()
		l.emit(tokenOperatorSub)
		return readExpr
	case unicode.IsSpace(r):
		return readDurationSpace
	}
	l.emit(tokenOperatorAdd)
	return readExpr
}

func readDurationSpace(l *lexer) stateFn {
	l.readFn(unicode.IsSpace)
	r := l.peek()
	switch {
	case r == eof:
		l.ignore()
		return nil
	case r == '+':
		fallthrough
	case r == '&':
		l.ignore()
		l.read()
		l.emit(tokenOperatorAdd)
		return readExpr
	case r == '-':
		l.ignore()
		l.read()
		l.emit(tokenOperatorSub)
		return readExpr
	case unicode.IsDigit(r):
		l.emit(tokenOperatorAdd)
		return readExpr
	}
	return readDurationSpaceNext
}

func readDurationSpaceNext(l *lexer) stateFn {
	space := l.value()
	l.ignore()
	v := l.peekFn(unicode.IsLetter)
	v = strings.ToLower(v)
	switch v {
	case "ago":
		l.readString("ago")
		l.emit(tokenAgo)
	case "before":
		l.readString("before")
		l.emit(tokenBefore)
	case "after":
		l.readString("after")
		l.emit(tokenFrom)
	case "from":
		l.readString("from")
		l.emit(tokenFrom)
	case "and":
		l.readString("and")
		l.emit(tokenOperatorAdd)
	case "one", "a":
		fallthrough
	case "two":
		fallthrough
	case "three":
		fallthrough
	case "four":
		fallthrough
	case "five":
		fallthrough
	case "six":
		fallthrough
	case "seven":
		fallthrough
	case "eight":
		fallthrough
	case "nine":
		fallthrough
	case "ten":
		fallthrough
	case "eleven":
		fallthrough
	case "twelve":
		l.emitAs(tokenOperatorAdd, space)
	}
	return readExpr
}

func readDurationUnitSecondsOrOrdinal(l *lexer) stateFn {
	l.read()
	r := l.peek()
	if r == 't' {
		l.read()
		l.emit(tokenOrdinal)
		return readExpr
	}
	l.emit(tokenUnit)
	return readDurationNextShort
}

func readLetter(l *lexer) stateFn {
	l.readFn(isTimeRune)
	v := l.value()
	v = strings.ToLower(v)
	switch v {
	case "now":
		l.emit(tokenNow)
		return readExpr
	case "today", "tomorrow", "yesterday":
		l.emit(tokenDate)
		return readExpr
	case "midnight", "noon":
		l.emit(tokenTime)
		return readExpr
	case "one", "a":
		l.emitAs(tokenDigit, "1")
		return readExpr
	case "two":
		l.emitAs(tokenDigit, "2")
		return readExpr
	case "three":
		l.emitAs(tokenDigit, "3")
		return readExpr
	case "four":
		l.emitAs(tokenDigit, "4")
		return readExpr
	case "five":
		l.emitAs(tokenDigit, "5")
		return readExpr
	case "six":
		l.emitAs(tokenDigit, "6")
		return readExpr
	case "seven":
		l.emitAs(tokenDigit, "7")
		return readExpr
	case "eight":
		l.emitAs(tokenDigit, "8")
		return readExpr
	case "nine":
		l.emitAs(tokenDigit, "9")
		return readExpr
	case "ten":
		l.emitAs(tokenDigit, "10")
		return readExpr
	case "eleven":
		l.emitAs(tokenDigit, "11")
		return readExpr
	case "twelve":
		l.emitAs(tokenDigit, "12")
		return readExpr
	case "am", "pm":
		l.emit(tokenTwelveHour)
		return readExpr
	case "year", "years":
		fallthrough
	case "month", "months":
		fallthrough
	case "week", "weeks":
		fallthrough
	case "day", "days":
		fallthrough
	case "hour", "hours":
		fallthrough
	case "minute", "minutes":
		fallthrough
	case "second", "seconds":
		l.emit(tokenUnit)
		return readDurationNext
	case "sun", "sunday":
		fallthrough
	case "mon", "monday":
		fallthrough
	case "tue", "tuesday":
		fallthrough
	case "wed", "wednesday":
		fallthrough
	case "thu", "thursday":
		fallthrough
	case "fri", "friday":
		fallthrough
	case "sat", "saturday":
		l.emit(tokenWeekday)
		return readExpr
	case "jan", "january":
		fallthrough
	case "feb", "february":
		fallthrough
	case "mar", "march":
		fallthrough
	case "apr", "april":
		fallthrough
	case "may":
		fallthrough
	case "jun", "june":
		fallthrough
	case "jul", "july":
		fallthrough
	case "aug", "august":
		fallthrough
	case "sep", "september":
		fallthrough
	case "oct", "october":
		fallthrough
	case "nov", "november":
		fallthrough
	case "dec", "december":
		l.emit(tokenMonth)
		return readExpr
	case "in", "of", "on", "the", "next", "last":
		fallthrough
	case "at", "quarter", "half", "past", "to", "after":
		fallthrough
	case "oclock", "o'clock", "morning", "afternoon", "evening":
		l.emit(tokenKeyword)
		return readExpr
	}
	return l.errorf("invalid character")
}

func readOrdinal(l *lexer) stateFn {
	var ok bool
	r := l.read()
	switch r {
	case 'n', 'r':
		ok = l.readRune('d') // 2nd, 3rd
	case 't':
		ok = l.readRune('h') // 4th-9th
	}
	if !ok {
		return l.errorf("invalid ordinal")
	}
	l.emit(tokenOrdinal)
	return readExpr
}

func readTwelveHour(l *lexer) stateFn {
	if !l.readString("am") && !l.readString("pm") {
		return l.errorf("expected twelve hour am/pm marker")
	}
	l.emit(tokenTwelveHour)
	return readExpr
}

func isTimeRune(r rune) bool {
	switch r {
	case '\'':
		return true
	}
	return unicode.IsLetter(r)
}
