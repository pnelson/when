// Package when implements a natural language date/time arithmetic parser.
package when

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	pos    int
	tokens []token
	now    time.Time
	lhs    []lhsFn
	rhs    time.Time
	sub    bool
	date   bool
	time   bool
}

// Parse returns the derived time.
func Parse(s string) (time.Time, error) {
	now := time.Now()
	return ParseNow(s, now)
}

// ParseNow returns the derived time relative to now.
func ParseNow(s string, now time.Time) (time.Time, error) {
	var t time.Time
	tokens, err := lex(s)
	if err != nil {
		return t, err
	}
	if len(tokens) == 0 {
		return now, nil
	}
	p := &parser{
		now:    now,
		tokens: tokens,
	}
	err = p.parseExpr()
	if err != nil {
		return t, err
	}
	t = p.rhs
	for _, fn := range p.lhs {
		t = fn.apply(t, p.sub)
	}
	return t, nil
}

func (p *parser) parseExpr() error {
	t := p.peek()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenDigit:
		return p.parseExprDigit()
	}
	return p.parseDateTime()
}

func (p *parser) parseExprDigit() error {
	d := p.next()
	t := p.peek()
	switch t.typ {
	case tokenEOF, tokenDateSeparator:
		return p.parseDateYear(d)
	case tokenUnit:
		return p.parseDurationLeftUnit(d, false)
	case tokenColon:
		return p.parseDigitColon(d)
	case tokenKeyword:
		return p.parseDigitKeyword(d)
	case tokenOrdinal:
		return p.parseDigitOrdinal(d)
	case tokenTwelveHour:
		t = p.next()
		m := strings.ToLower(t.val)
		return p.parseDigitTwelveHour(d, m)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDateTime() error {
	t := p.peek()
	switch t.typ {
	case tokenNow:
		p.next()
		return p.parseNow()
	case tokenDate:
		return p.parseDateConst()
	case tokenMonth:
		return p.parseMonth()
	case tokenWeekday:
		return p.parseWeekday()
	case tokenKeyword:
		return p.parseKeyword()
	case tokenTime:
		return p.parseTimeConst()
	case tokenDigit:
		t = p.next()
		return p.parseDigit(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDate() error {
	p.time = true
	if p.date {
		return p.parseDurationRightNext()
	}
	t := p.peek()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenOperatorAdd:
		return p.parseDurationRightNext()
	case tokenOperatorSub:
		return p.parseDurationRightNext()
	case tokenDate:
		return p.parseDateConst()
	case tokenMonth:
		return p.parseMonth()
	case tokenWeekday:
		return p.parseWeekday()
	case tokenKeyword:
		return p.parseDateKeyword()
	case tokenDigit:
		t = p.next()
		return p.parseDigit(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDateConst() error {
	t := p.next()
	loc := p.now.Location()
	y, M, d := p.now.Date()
	h, m, s := p.rhs.Clock()
	switch t.val {
	case "today":
	case "tomorrow":
		d++
	case "yesterday":
		d--
	default:
		return newParseError(t, "unexpected date")
	}
	p.rhs = time.Date(y, M, d, h, m, s, 0, loc)
	return p.parseTime()
}

func (p *parser) parseDateKeyword() error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "on" {
		return newParseError(t, "unexpected token")
	}
	return p.parseKeywordOn()
}

func (p *parser) parseDateYear(d token) error {
	y, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(y, time.January, 1, h, m, s, 0, loc)
	return p.parseDateYearMonth()
}

func (p *parser) parseDateYearMonth() error {
	t := p.peek()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenDateSeparator:
		p.next()
	default:
		return p.parseTime()
	}
	t = p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	M, err := strconv.Atoi(t.val)
	if err != nil {
		return err
	}
	loc := p.rhs.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.rhs.Year(), time.Month(M), 1, h, m, s, 0, loc)
	return p.parseDateYearMonthDay()
}

func (p *parser) parseDateYearMonthDay() error {
	t := p.peek()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenDateSeparator:
		p.next()
	default:
		return p.parseTime()
	}
	t = p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	d, err := strconv.Atoi(t.val)
	if err != nil {
		return err
	}
	loc := p.rhs.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.rhs.Year(), p.rhs.Month(), d, h, m, s, 0, loc)
	return p.parseTime()
}

func (p *parser) parseDigit(d token) error {
	t := p.peek()
	switch t.typ {
	case tokenEOF, tokenDateSeparator:
		return p.parseDateYear(d)
	case tokenColon:
		return p.parseDigitColon(d)
	case tokenKeyword:
		return p.parseDigitKeyword(d)
	case tokenOrdinal:
		return p.parseDigitOrdinal(d)
	case tokenTwelveHour:
		t = p.next()
		m := strings.ToLower(t.val)
		return p.parseDigitTwelveHour(d, m)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitColon(h token) error {
	p.next()
	m := p.next()
	if m.typ != tokenDigit {
		return newParseError(m, "unexpected token")
	}
	t := p.peek()
	if t.typ == tokenTwelveHour {
		return p.parseDigitColonTwelveHour(h, m)
	}
	return p.parseDigitColonTwentyFourHour(h, m)
}

func (p *parser) parseDigitColonTwelveHour(h, m token) error {
	t := p.next()
	i := strings.ToLower(t.val)
	loc := p.now.Location()
	r, err := time.ParseInLocation("3:04pm", h.val+":"+m.val+i, loc)
	if err != nil {
		return err
	}
	y, M, d := p.rhs.Date()
	if p.rhs.IsZero() {
		y, M, d = p.now.Date()
	}
	p.rhs = time.Date(y, M, d, r.Hour(), r.Minute(), 0, 0, loc)
	return p.parseDate()
}

func (p *parser) parseDigitColonTwentyFourHour(h, m token) error {
	t := p.peek()
	if t.typ == tokenColon {
		p.next()
		return p.parseDigitColonTwentyFourHourWithSeconds(h, m)
	}
	loc := p.now.Location()
	r, err := time.ParseInLocation("15:04", h.val+":"+m.val, loc)
	if err != nil {
		return err
	}
	y, M, d := p.rhs.Date()
	if p.rhs.IsZero() {
		y, M, d = p.now.Date()
	}
	p.rhs = time.Date(y, M, d, r.Hour(), r.Minute(), 0, 0, loc)
	return p.parseDate()
}

func (p *parser) parseDigitColonTwentyFourHourWithSeconds(h, m token) error {
	s := p.next()
	if s.typ != tokenDigit {
		return newParseError(m, "unexpected token")
	}
	loc := p.now.Location()
	r, err := time.ParseInLocation("15:04:05", h.val+":"+m.val+":"+s.val, loc)
	if err != nil {
		return err
	}
	y, M, d := p.rhs.Date()
	if p.rhs.IsZero() {
		y, M, d = p.now.Date()
	}
	p.rhs = time.Date(y, M, d, r.Hour(), r.Minute(), r.Second(), 0, loc)
	return p.parseDate()
}

func (p *parser) parseDigitKeyword(d token) error {
	t := p.next()
	if t.typ != tokenKeyword {
		return newParseError(t, "unexpected token")
	}
	switch t.val {
	case "in":
		return p.parseDigitKeywordIn(d)
	case "oclock", "o'clock":
		return p.parseDigitKeywordOclock(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitKeywordIn(d token) error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "the" {
		return newParseError(t, "unexpected token")
	}
	t = p.next()
	if t.typ != tokenKeyword {
		return newParseError(t, "unexpected token")
	}
	switch t.val {
	case "morning":
		return p.parseDigitTwelveHour(d, "am")
	case "afternoon", "evening":
		return p.parseDigitTwelveHour(d, "pm")
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitKeywordOclock(d token) error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "in" {
		return newParseError(t, "unexpected token")
	}
	return p.parseDigitKeywordIn(d)
}

func (p *parser) parseDigitOrdinal(d token) error {
	t := p.next()
	if t.typ != tokenOrdinal {
		return newParseError(t, "unexpected token")
	}
	n, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	t = p.peek()
	switch t.typ {
	case tokenEOF:
		return p.parseDigitOrdinalEOF(n)
	case tokenKeyword:
		return p.parseDigitOrdinalKeyword(n)
	case tokenWeekday:
		return p.parseDigitOrdinalWeekday(n)
	case tokenMonth:
		return p.parseDigitOrdinalMonth(n)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalEOF(d int) error {
	loc := p.now.Location()
	y, M, _ := p.now.Date()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(y, M, d, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(0, 1, 0)
	}
	return nil
}

func (p *parser) parseDigitOrdinalKeyword(d int) error {
	t := p.next()
	if t.typ != tokenKeyword {
		return newParseError(t, "unexpected token")
	}
	switch t.val {
	case "@", "at":
		return p.parseDigitOrdinalAt(d)
	case "of":
		return p.parseDigitOrdinalOf(d)
	case "last":
		return p.parseDigitOrdinalLast(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalAt(d int) error {
	loc := p.now.Location()
	y, M, _ := p.now.Date()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(y, M, d, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(0, 1, 0)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalOf(d int) error {
	t := p.peek()
	switch t.typ {
	case tokenKeyword:
		return p.parseDigitOrdinalOfKeyword(d)
	case tokenMonth:
		return p.parseDigitOrdinalOfMonth(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalOfKeyword(d int) error {
	t := p.next()
	u := p.next()
	if u.typ != tokenUnit || u.val != "month" {
		return newParseError(u, "unexpected token")
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), p.now.Month(), d, h, m, s, 0, loc)
	switch t.val {
	case "the":
	case "last":
		p.rhs = p.rhs.AddDate(0, -1, 0)
	case "next":
		p.rhs = p.rhs.AddDate(0, 1, 0)
	default:
		return newParseError(t, "unexpected token")
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalOfMonth(d int) error {
	t := p.next()
	M, err := parseMonth(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, d, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalLast(d int) error {
	t := p.peek()
	switch t.typ {
	case tokenUnit:
		return p.parseDigitOrdinalLastDay(d)
	case tokenWeekday:
		return p.parseDigitOrdinalLastWeekday(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalLastDay(d int) error {
	t := p.next()
	if t.typ != tokenUnit || t.val != "day" {
		return newParseError(t, "unexpected token")
	}
	return p.parseDigitOrdinalLastDayOf(d)
}

func (p *parser) parseDigitOrdinalLastDayOf(d int) error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "of" && t.val != "in" {
		return newParseError(t, "unexpected token")
	}
	t = p.peek()
	switch t.typ {
	case tokenKeyword:
		return p.parseDigitOrdinalLastDayOfKeyword(d)
	case tokenMonth:
		return p.parseDigitOrdinalLastDayOfMonth(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalLastDayOfKeyword(d int) error {
	t := p.next()
	u := p.next()
	if u.typ != tokenUnit || u.val != "month" {
		return newParseError(u, "unexpected token")
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), p.now.Month(), 1, h, m, s, 0, loc)
	switch t.val {
	case "the":
		p.rhs = p.rhs.AddDate(0, 1, -d)
	case "last":
		p.rhs = p.rhs.AddDate(0, 0, -d)
	case "next":
		p.rhs = p.rhs.AddDate(0, 2, -d)
	default:
		return newParseError(t, "unexpected token")
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalLastDayOfMonth(d int) error {
	t := p.next()
	M, err := parseMonth(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, 1, h, m, s, 0, loc)
	p.rhs = p.rhs.AddDate(0, 1, -d)
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalLastWeekday(d int) error {
	t := p.next()
	w, err := parseWeekday(t)
	if err != nil {
		return err
	}
	t = p.peek()
	if t.typ == tokenKeyword && (t.val == "of" || t.val == "in") {
		p.next()
		return p.parseDigitOrdinalLastWeekdayOf(d, w)
	}
	loc := p.now.Location()
	y, M, d := p.now.Date()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(y, M, d, h, m, s, 0, loc)
	p.rhs = p.rhs.AddDate(0, 0, int(w-p.rhs.Weekday()-7))
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalLastWeekdayOf(d int, w time.Weekday) error {
	t := p.peek()
	switch t.typ {
	case tokenKeyword:
		return p.parseDigitOrdinalLastWeekdayOfKeyword(d, w)
	case tokenMonth:
		return p.parseDigitOrdinalLastWeekdayOfMonth(d, w)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalLastWeekdayOfKeyword(d int, w time.Weekday) error {
	t := p.next()
	u := p.next()
	if u.typ != tokenUnit || u.val != "month" {
		return newParseError(t, "unexpected token")
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), p.now.Month(), 1, h, m, s, 0, loc)
	switch t.val {
	case "the":
		p.rhs = p.rhs.AddDate(0, 1, -1)
	case "last":
		p.rhs = p.rhs.AddDate(0, 0, -1)
	case "next":
		p.rhs = p.rhs.AddDate(0, 2, -1)
	default:
		return newParseError(t, "unexpected token")
	}
	days := int(p.rhs.Weekday() - w)
	if days < 0 {
		days += 7
	}
	p.rhs = p.rhs.AddDate(0, 0, -days)
	for i := 0; i < d-1; i++ {
		p.rhs = p.rhs.AddDate(0, 0, -7)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalLastWeekdayOfMonth(d int, w time.Weekday) error {
	t := p.next()
	M, err := parseMonth(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, 1, h, m, s, 0, loc)
	p.rhs = p.rhs.AddDate(0, 1, -1)
	days := int(p.rhs.Weekday() - w)
	if days < 0 {
		days += 7
	}
	p.rhs = p.rhs.AddDate(0, 0, -days)
	for i := 0; i < d-1; i++ {
		p.rhs = p.rhs.AddDate(0, 0, -7)
	}
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalWeekday(d int) error {
	t := p.next()
	w, err := parseWeekday(t)
	if err != nil {
		return err
	}
	return p.parseDigitOrdinalWeekdayOf(d, w)
}

func (p *parser) parseDigitOrdinalWeekdayOf(d int, w time.Weekday) error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "of" && t.val != "in" {
		return newParseError(t, "unexpected token")
	}
	t = p.peek()
	switch t.typ {
	case tokenKeyword:
		return p.parseDigitOrdinalWeekdayOfKeyword(d, w)
	case tokenMonth:
		return p.parseDigitOrdinalWeekdayOfMonth(d, w)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigitOrdinalWeekdayOfKeyword(d int, w time.Weekday) error {
	t := p.next()
	u := p.next()
	if u.typ != tokenUnit || u.val != "month" {
		return newParseError(u, "unexpected token")
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), p.now.Month(), 1, h, m, s, 0, loc)
	switch t.val {
	case "the":
	case "last":
		p.rhs = p.rhs.AddDate(0, -1, 0)
	case "next":
		p.rhs = p.rhs.AddDate(0, 1, 0)
	default:
		return newParseError(t, "unexpected token")
	}
	days := int(w - p.rhs.Weekday())
	if days < 0 {
		days += 7
	}
	p.rhs = p.rhs.AddDate(0, 0, days)
	for i := 0; i < d-1; i++ {
		p.rhs = p.rhs.AddDate(0, 0, 7)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalWeekdayOfMonth(d int, w time.Weekday) error {
	t := p.next()
	M, err := parseMonth(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, 1, h, m, s, 0, loc)
	days := int(w - p.rhs.Weekday())
	if days < 0 {
		days += 7
	}
	p.rhs = p.rhs.AddDate(0, 0, days)
	for i := 0; i < d-1; i++ {
		p.rhs = p.rhs.AddDate(0, 0, 7)
	}
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return p.parseTime()
}

func (p *parser) parseDigitOrdinalMonth(d int) error {
	t := p.next()
	M, err := parseMonth(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, d, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return p.parseTime()
}

func (p *parser) parseDigitTwelveHour(h token, i string) error {
	r, err := time.ParseInLocation("3pm", h.val+i, p.now.Location())
	if err != nil {
		return err
	}
	y, M, d := p.rhs.Date()
	if p.rhs.IsZero() {
		y, M, d = p.now.Date()
	}
	p.rhs = time.Date(y, M, d, r.Hour(), 0, 0, 0, r.Location())
	return p.parseDate()
}

func (p *parser) parseDurationLeft(sub bool) error {
	t := p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	return p.parseDurationLeftUnit(t, sub)
}

func (p *parser) parseDurationLeftAgo() error {
	t := p.next()
	if t.typ != tokenEOF {
		return newParseError(t, "unexpected token")
	}
	p.sub = true
	p.rhs = p.now
	return nil
}

func (p *parser) parseDurationLeftBefore() error {
	t := p.peek()
	if t.typ == tokenEOF {
		return newParseError(t, "unexpected token")
	}
	p.sub = true
	return p.parseDateTime()
}

func (p *parser) parseDurationLeftFrom() error {
	t := p.peek()
	if t.typ == tokenEOF {
		return newParseError(t, "unexpected token")
	}
	return p.parseDateTime()
}

func (p *parser) parseDurationLeftNext() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.rhs = p.now
		return nil
	case tokenOperatorAdd:
		return p.parseDurationLeft(false)
	case tokenOperatorSub:
		return p.parseDurationLeft(true)
	case tokenAgo:
		return p.parseDurationLeftAgo()
	case tokenBefore:
		return p.parseDurationLeftBefore()
	case tokenFrom:
		return p.parseDurationLeftFrom()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDurationLeftUnit(d token, sub bool) error {
	u := p.next()
	n, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	if sub {
		n *= -1
	}
	p.lhs = append(p.lhs, lhsFn{n, u.val})
	return p.parseDurationLeftNext()
}

func (p *parser) parseDurationRight(sub bool) error {
	t := p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	return p.parseDurationRightUnit(t, sub)
}

func (p *parser) parseDurationRightNext() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenOperatorAdd:
		return p.parseDurationRight(false)
	case tokenOperatorSub:
		return p.parseDurationRight(true)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDurationRightUnit(d token, sub bool) error {
	u := p.next()
	n, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	if sub {
		n *= -1
	}
	switch u.val {
	case "y", "year", "years":
		p.rhs = p.rhs.AddDate(n, 0, 0)
	case "M", "month", "months":
		p.rhs = p.rhs.AddDate(0, n, 0)
	case "w", "week", "weeks":
		p.rhs = p.rhs.AddDate(0, 0, 7*n)
	case "d", "day", "days":
		p.rhs = p.rhs.AddDate(0, 0, n)
	case "h", "hour", "hours":
		p.rhs = p.rhs.Add(time.Duration(n) * time.Hour)
	case "m", "minute", "minutes":
		p.rhs = p.rhs.Add(time.Duration(n) * time.Minute)
	case "s", "second", "seconds":
		p.rhs = p.rhs.Add(time.Duration(n) * time.Second)
	default:
		return newParseError(u, "unexpected token")
	}
	return p.parseDurationRightNext()
}

func (p *parser) parseKeyword() error {
	t := p.next()
	if t.typ != tokenKeyword {
		return newParseError(t, "unexpected token")
	}
	switch t.val {
	case "@", "at":
		return p.parseKeywordAt()
	case "on":
		return p.parseKeywordOn()
	case "last":
		return p.parseDigitOrdinalLast(1)
	case "half":
		return p.parseKeywordHalf()
	case "quarter":
		return p.parseKeywordQuarter()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseKeywordAt() error {
	t := p.peek()
	switch t.typ {
	case tokenTime:
		return p.parseTimeConst()
	case tokenDigit:
		t = p.next()
		return p.parseDigit(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseKeywordHalf() error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "past" {
		return newParseError(t, "unexpected token")
	}
	return p.parseKeywordHalfPast()
}

func (p *parser) parseKeywordHalfPast() error {
	t := p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	err := p.parseDigit(t)
	if err != nil {
		return err
	}
	p.rhs = p.rhs.Add(30 * time.Minute)
	return nil
}

func (p *parser) parseKeywordOn() error {
	t := p.peek()
	switch t.typ {
	case tokenDigit:
		d := p.next()
		return p.parseDateYear(d)
	case tokenWeekday:
		return p.parseWeekday()
	case tokenMonth:
		return p.parseMonth()
	case tokenKeyword:
		return p.parseKeywordOnThe()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseKeywordOnThe() error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "the" {
		return newParseError(t, "unexpected token")
	}
	t = p.peek()
	switch t.typ {
	case tokenDigit:
		return p.parseKeywordOnTheDigit()
	case tokenKeyword:
		return p.parseKeywordOnTheLast()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseKeywordOnTheLast() error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "last" {
		return newParseError(t, "unexpected token")
	}
	return p.parseDigitOrdinalLast(1)
}

func (p *parser) parseKeywordOnTheDigit() error {
	d := p.next()
	if d.typ != tokenDigit {
		return newParseError(d, "unexpected token")
	}
	t := p.peek()
	if t.typ != tokenOrdinal {
		return newParseError(t, "unexpected token")
	}
	return p.parseDigitOrdinal(d)
}

func (p *parser) parseKeywordQuarter() error {
	t := p.next()
	switch t.val {
	case "to":
		return p.parseKeywordQuarterTo()
	case "after", "past":
		return p.parseKeywordQuarterAfter()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseKeywordQuarterTo() error {
	t := p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	err := p.parseDigit(t)
	if err != nil {
		return err
	}
	p.rhs = p.rhs.Add(-15 * time.Minute)
	return nil
}

func (p *parser) parseKeywordQuarterAfter() error {
	t := p.next()
	if t.typ != tokenDigit {
		return newParseError(t, "unexpected token")
	}
	err := p.parseDigit(t)
	if err != nil {
		return err
	}
	p.rhs = p.rhs.Add(15 * time.Minute)
	return nil
}

func (p *parser) parseMonth() error {
	t := p.next()
	m, err := parseMonth(t)
	if err != nil {
		return err
	}
	t = p.peek()
	switch t.typ {
	case tokenEOF:
		return p.parseMonthEOF(m)
	case tokenKeyword:
		return p.parseMonthThe(m)
	case tokenDigit:
		return p.parseMonthTheDigit(m)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseMonthEOF(M time.Month) error {
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, 1, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return nil
}

func (p *parser) parseMonthThe(M time.Month) error {
	t := p.next()
	if t.typ != tokenKeyword || t.val != "the" {
		return newParseError(t, "unexpected token")
	}
	return p.parseMonthTheDigit(M)
}

func (p *parser) parseMonthTheDigit(M time.Month) error {
	d := p.next()
	if d.typ != tokenDigit {
		return newParseError(d, "unexpected token")
	}
	t := p.next()
	if t.typ != tokenOrdinal {
		return newParseError(t, "unexpected token")
	}
	n, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	h, m, s := p.rhs.Clock()
	p.rhs = time.Date(p.now.Year(), M, n, h, m, s, 0, loc)
	if p.now.Equal(p.rhs) || p.now.After(p.rhs) {
		p.rhs = p.rhs.AddDate(1, 0, 0)
	}
	return p.parseTime()
}

func (p *parser) parseNow() error {
	p.rhs = p.now
	t := p.next()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenOperatorAdd:
		return p.parseDurationRight(false)
	case tokenOperatorSub:
		return p.parseDurationRight(true)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseTime() error {
	p.date = true
	if p.time {
		return p.parseDurationRightNext()
	}
	t := p.peek()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenOperatorAdd:
		return p.parseDurationRightNext()
	case tokenOperatorSub:
		return p.parseDurationRightNext()
	case tokenTime:
		return p.parseTimeConst()
	case tokenKeyword:
		return p.parseTimeKeyword()
	case tokenDigit:
		t = p.next()
		return p.parseDigit(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseTimeConst() error {
	t := p.next()
	loc := p.now.Location()
	y, M, d := p.rhs.Date()
	if p.rhs.IsZero() {
		y, M, d = p.now.Date()
	}
	switch t.val {
	case "midnight":
		p.rhs = time.Date(y, M, d, 0, 0, 0, 0, loc)
	case "noon":
		p.rhs = time.Date(y, M, d, 12, 0, 0, 0, loc)
	default:
		return newParseError(t, "unexpected date")
	}
	return p.parseDate()
}

func (p *parser) parseTimeKeyword() error {
	t := p.next()
	if t.typ != tokenKeyword {
		return newParseError(t, "unexpected token")
	}
	switch t.val {
	case "@", "at":
		return p.parseKeywordAt()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseWeekday() error {
	t := p.next()
	w, err := parseWeekday(t)
	if err != nil {
		return err
	}
	loc := p.now.Location()
	y, M, d := p.now.Date()
	h, m, s := p.rhs.Clock()
	days := int(w - p.now.Weekday())
	if days <= 0 {
		days += 7
	}
	p.rhs = time.Date(y, M, d+days, h, m, s, 0, loc)
	return p.parseTime()
}

func (p *parser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{tokenEOF, ""}
	}
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	t := p.peek()
	p.pos++
	return t
}

func parseMonth(t token) (time.Month, error) {
	var m time.Month
	if t.typ != tokenMonth {
		return m, newParseError(t, "unexpected token")
	}
	switch strings.ToLower(t.val[:3]) {
	case "jan":
		m = time.January
	case "feb":
		m = time.February
	case "mar":
		m = time.March
	case "apr":
		m = time.April
	case "may":
		m = time.May
	case "jun":
		m = time.June
	case "jul":
		m = time.July
	case "aug":
		m = time.August
	case "sep":
		m = time.September
	case "oct":
		m = time.October
	case "nov":
		m = time.November
	case "dec":
		m = time.December
	default:
		return m, newParseError(t, "invalid month")
	}
	return m, nil
}

func parseWeekday(t token) (time.Weekday, error) {
	var w time.Weekday
	if t.typ != tokenWeekday {
		return w, newParseError(t, "unexpected token")
	}
	switch strings.ToLower(t.val[:3]) {
	case "sun":
		w = time.Sunday
	case "mon":
		w = time.Monday
	case "tue":
		w = time.Tuesday
	case "wed":
		w = time.Wednesday
	case "thu":
		w = time.Thursday
	case "fri":
		w = time.Friday
	case "sat":
		w = time.Saturday
	default:
		return w, newParseError(t, "invalid weekday")
	}
	return w, nil
}

type lhsFn struct {
	n    int
	unit string
}

func (f lhsFn) apply(t time.Time, sub bool) time.Time {
	if sub {
		f.n *= -1
	}
	switch f.unit {
	case "y", "year", "years":
		return t.AddDate(f.n, 0, 0)
	case "M", "month", "months":
		return t.AddDate(0, f.n, 0)
	case "w", "week", "weeks":
		return t.AddDate(0, 0, 7*f.n)
	case "d", "day", "days":
		return t.AddDate(0, 0, f.n)
	case "h", "hour", "hours":
		return t.Add(time.Duration(f.n) * time.Hour)
	case "m", "minute", "minutes":
		return t.Add(time.Duration(f.n) * time.Minute)
	case "s", "second", "seconds":
		return t.Add(time.Duration(f.n) * time.Second)
	}
	return time.Time{}
}

type parseError struct {
	token   token
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("%s, token: %q", e.message, e.token.val)
}

func newParseError(t token, message string) parseError {
	return parseError{t, message}
}
