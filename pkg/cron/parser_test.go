package cron_test

import (
	"cron/pkg/cron"
	"errors"
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
	{"* -2 * * *", cron.ValueTooLow},
	{"* 24 * * *", cron.ValueTooHigh},
	{"* * 0 * *", cron.ValueTooLow},
	{"* * 32 * *", cron.ValueTooHigh},
	{"* * * 0 *", cron.ValueTooLow},
	{"* * * 13 *", cron.ValueTooHigh},
	{"* * * * -1", cron.ValueTooLow},
	{"* * * * 7", cron.ValueTooHigh},
}

func TestInvalidPatterns(t *testing.T) {
	for _, tc := range invalidPatterns {
		_, err := cron.Parse(tc.pattern)
		if !errors.Is(err, tc.err) {
			t.Errorf("the following error: \"%v\" is expected for a pattern \"%s\". Got \"%v\"", tc.err, tc.pattern, err)
		}
	}
}
