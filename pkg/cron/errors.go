package cron

import "errors"

var (
	FieldsCountInvalid        = errors.New("an incorrect number of fields in the cron config")
	PatternFieldsCountInvalid = errors.New("an incorrect number of fields in the cron pattern")
	NumericalValueInvalid     = errors.New("number value invalid")
	ValueTooLow               = errors.New("value lower than minimum")
	ValueTooHigh              = errors.New("value higher than maximum")
	StepPatternIncomplete     = errors.New("part of step missing")
	StepPatternValueInvalid   = errors.New("value invalid for a step")
	StepPatternValueTooLow    = errors.New("value for a step lower than minimum")
	StepPatternValueTooHigh   = errors.New("value for a step higher than maximum")
	StepValueInvalid          = errors.New("step value incorrect")
	StepValueOutsideRange     = errors.New("step value outside valid range")
	RangePatternBoundTooLow   = errors.New("lower bound value lower than minimum")
	RangePatternBoundTooHigh  = errors.New("upper bound value higher than maximum")
)
