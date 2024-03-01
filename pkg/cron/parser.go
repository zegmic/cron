package cron

import (
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Minutes []byte
	Hours   []byte
	Days    []byte
	Months  []byte
	WeekDay []byte
	Command string
}

func Parse(pattern string) (*Config, error) {
	params := strings.Split(pattern, " ")
	if len(params) != 6 {
		return nil, FieldsCountInvalid
	}

	mins, err := convert(params[0], 0, 59)
	if err != nil {
		return nil, err
	}
	hours, err := convert(params[1], 0, 23)
	if err != nil {
		return nil, err
	}
	days, err := convert(params[2], 1, 31)
	if err != nil {
		return nil, err
	}
	months, err := convert(params[3], 1, 12)
	if err != nil {
		return nil, err
	}
	weekDay, err := convert(params[4], 0, 6)
	if err != nil {
		return nil, err
	}
	return &Config{
		Minutes: mins,
		Hours:   hours,
		Days:    days,
		Months:  months,
		WeekDay: weekDay,
		Command: params[5],
	}, nil
}

func convert(field string, min, max byte) ([]byte, error) {
	if strings.Contains(field, ",") {
		return convertList(field, min, max)
	} else if strings.Contains(field, "-") {
		return convertRange(field, min, max)
	} else if strings.Contains(field, "/") {
		return convertStep(field, min, max)
	} else {
		return convertValue(field, min, max)
	}
}

func convertValue(field string, min, max byte) ([]byte, error) {
	if field == "*" {
		var res []byte
		for i := min; i <= max; i++ {
			res = append(res, i)
		}
		return res, nil
	}

	val, err := strconv.Atoi(field)
	if err != nil {
		return nil, NumericalValueInvalid
	}
	if val < int(min) {
		return nil, ValueTooLow
	}
	if byte(val) > max {
		return nil, ValueTooHigh
	}

	return []byte{byte(val)}, nil
}

func convertStep(field string, min, max byte) ([]byte, error) {
	stepPattern := strings.Split(field, "/")
	if len(stepPattern) != 2 {
		return nil, StepPatternIncomplete
	}

	var start byte
	if stepPattern[0] == "*" {
		start = 0
	} else {
		val, err := strconv.Atoi(stepPattern[0])
		if err != nil {
			return nil, StepPatternValueInvalid
		}
		if byte(val) < min {
			return nil, StepPatternValueTooLow
		}
		start = byte(val)
	}

	step, err := strconv.Atoi(stepPattern[1])
	if err != nil {
		return nil, StepValueInvalid
	}

	if step <= 0 || byte(step) > max {
		return nil, StepValueOutsideRange
	}

	var res []byte
	for i := start; i <= max; i += byte(step) {
		res = append(res, i)
	}

	return res, nil
}

func convertRange(field string, min byte, max byte) ([]byte, error) {
	minmax := strings.Split(field, "-")
	if minmax[0] == "" {
		return convertValue(field, min, max)
	}
	minVal, err := strconv.Atoi(minmax[0])
	if err != nil {
		return nil, NumericalValueInvalid
	}
	if byte(minVal) < min {
		return nil, RangePatternBoundTooLow
	}

	maxVal, err := strconv.Atoi(minmax[1])
	if err != nil {
		return nil, NumericalValueInvalid
	}
	if byte(maxVal) > max {
		return nil, RangePatternBoundTooHigh
	}

	var res []byte
	for i := byte(minVal); i <= byte(maxVal); i++ {
		res = append(res, i)
	}

	return res, nil
}

func convertList(field string, min byte, max byte) ([]byte, error) {
	subfields := strings.Split(field, ",")
	var res []byte
	for _, sub := range subfields {
		converted, err := convert(sub, min, max)
		if err != nil {
			return nil, err
		}
		res = append(res, converted...)
	}

	return removeDups(res), nil
}

func removeDups(res []byte) []byte {
	var unique []byte
	keys := make(map[byte]bool)
	for _, v := range res {
		if !keys[v] {
			keys[v] = true
			unique = append(unique, v)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i] < unique[j]
	})

	return unique
}
