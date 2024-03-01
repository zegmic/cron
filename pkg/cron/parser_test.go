package cron

import (
	"errors"
	"testing"
)

var invalidPatterns = []struct {
	pattern string
	err     error
}{
	{"", FieldsCountInvalid},
	{"2 4 1 4 5", FieldsCountInvalid},
	{"2 4 1 4 5 6 4", FieldsCountInvalid},
	{"-12 * * * * *", ValueTooLow},
	{"60 * * * * *", ValueTooHigh},
}

func TestInvalidPatterns(t *testing.T) {
	for _, tc := range invalidPatterns {
		_, err := Parse(tc.pattern)
		if !errors.Is(err, tc.err) {
			t.Errorf("the following error: \"%v\" is expected for a pattern \"%s\". Got \"%v\"", tc.err, tc.pattern, err)
		}
	}
}
