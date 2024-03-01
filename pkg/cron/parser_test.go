package cron_test

import (
	"cron/pkg/cron"
	"errors"
	"reflect"
	"testing"
)

var invalidPatterns = []struct {
	pattern string
	err     error
}{
	{"", cron.PatternFieldsCountInvalid},
	{"2 4 1 4", cron.PatternFieldsCountInvalid},
	{"2 4 1 4 5 6 4", cron.PatternFieldsCountInvalid},
	{"-12 * * * *", cron.ValueTooLow},
	{"60 * * * *", cron.ValueTooHigh},
	{"a * * * *", cron.NumericalValueInvalid},
	{"* -2 * * *", cron.ValueTooLow},
	{"* 24 * * *", cron.ValueTooHigh},
	{"* a * * *", cron.NumericalValueInvalid},
	{"* * 0 * *", cron.ValueTooLow},
	{"* * 32 * *", cron.ValueTooHigh},
	{"* * a * *", cron.NumericalValueInvalid},
	{"* * * 0 *", cron.ValueTooLow},
	{"* * * 13 *", cron.ValueTooHigh},
	{"* * * a *", cron.NumericalValueInvalid},
	{"* * * * -1", cron.ValueTooLow},
	{"* * * * 7", cron.ValueTooHigh},
	{"* * * * a", cron.NumericalValueInvalid},
	{"-2-55 * * * *", cron.RangePatternBoundTooLow},
	{"6-60 * * * *", cron.RangePatternBoundTooHigh},
	{"* -5-10 * * *", cron.RangePatternBoundTooLow},
	{"* 2-25 * * *", cron.RangePatternBoundTooHigh},
	{"* * -3-10 * *", cron.RangePatternBoundTooLow},
	{"* * 2-32 * *", cron.RangePatternBoundTooHigh},
	{"* * * -6-8 *", cron.RangePatternBoundTooLow},
	{"* * * 4-13 *", cron.RangePatternBoundTooHigh},
	{"* * * * -1-3", cron.RangePatternBoundTooLow},
	{"* * * * 2-7", cron.RangePatternBoundTooHigh},
	{"/3 * * * *", cron.StepPatternIncomplete},
	{"5/ * * * *", cron.StepPatternIncomplete},
	{"/ * * * *", cron.StepPatternIncomplete},
	{"-5/4 * * * *", cron.StepPatternValueTooLow},
	{"123/4 * * * *", cron.StepPatternValueTooHigh},
	{"abc/4 * * * *", cron.StepPatternValueInvalid},
	{"4/-12 * * * *", cron.StepValueOutsideRange},
	{"6/60 * * * *", cron.StepValueOutsideRange},
	{"-5-8,-10-20,-15-30 0 1 1 0", cron.RangePatternBoundTooLow},
}

func TestInvalidPatterns(t *testing.T) {
	for _, tc := range invalidPatterns {
		_, err := cron.Parse(tc.pattern)
		if !errors.Is(err, tc.err) {
			t.Errorf("the following error: \"%v\" is expected for a pattern \"%s\". Got \"%v\"", tc.err, tc.pattern, err)
		}
	}
}

var validPatterns = []struct {
	pattern  string
	expected *cron.Pattern
}{
	{"1 0 1 1 0", newPattern([]byte{1}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"5-10 0 1 1 0", newPattern([]byte{5, 6, 7, 8, 9, 10}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"3/11 0 1 1 0", newPattern([]byte{3, 14, 25, 36, 47, 58}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"*/18 0 1 1 0", newPattern([]byte{0, 18, 36, 54}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"5,10,15 0 1 1 0", newPattern([]byte{5, 10, 15}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"5-8,10,15 0 1 1 0", newPattern([]byte{5, 6, 7, 8, 10, 15}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"*/15,7,13-16 0 1 1 0", newPattern([]byte{0, 7, 13, 14, 15, 16, 30, 45}, []byte{0}, []byte{1}, []byte{1}, []byte{0})},
	{"*/15,7,13-16 0 1 * 0", newPattern([]byte{0, 7, 13, 14, 15, 16, 30, 45}, []byte{0}, []byte{1}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, []byte{0})},
	{"*/15 0 1,15 * 1-5", newPattern([]byte{0, 15, 30, 45}, []byte{0}, []byte{1, 15}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, []byte{1, 2, 3, 4, 5})},
	{"* * * * *", newPattern(newRange(0, 59), newRange(0, 23), newRange(1, 31), newRange(1, 12), newRange(0, 6))},
}

func newPattern(mins, hours, days, months, weekDays []byte) *cron.Pattern {
	return &cron.Pattern{
		Minutes: mins,
		Hours:   hours,
		Days:    days,
		Months:  months,
		WeekDay: weekDays,
	}
}

func newRange(min, max byte) []byte {
	var res []byte
	for i := min; i <= max; i++ {
		res = append(res, i)
	}
	return res
}

func TestValidPatterns(t *testing.T) {
	for _, tc := range validPatterns {
		p, err := cron.Parse(tc.pattern)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(p, tc.expected) {
			t.Errorf("expected output %v for pattern %s but got %v", tc.expected, tc.pattern, p)
		}
	}
}

func TestInvalidConfig(t *testing.T) {
	args := "*/15 0 1,15 * 1-5"
	_, err := cron.ParseConfig(args)
	if !errors.Is(err, cron.FieldsCountInvalid) {
		t.Errorf("expected the following error: %v got: %v", cron.FieldsCountInvalid, err)
	}
}

func TestParseConfig(t *testing.T) {
	args := "*/15 0 1,15 * 1-5 /usr/bin/find"
	config, err := cron.ParseConfig(args)
	if err != nil {
		t.Error(err)
	}
	expected := &cron.Config{
		Pattern: cron.Pattern{
			Minutes: []byte{0, 15, 30, 45},
			Hours:   []byte{0},
			Days:    []byte{1, 15},
			Months:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			WeekDay: []byte{1, 2, 3, 4, 5},
		},
		Command: "/usr/bin/find",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("expected %v for arguments %s but got %v", expected, args, config)
	}
}
