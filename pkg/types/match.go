package types

import (
	"regexp"
	"strings"
	"time"
)

func stringMatch(op string, qValue string, oValue string) bool {
	test := regexp.QuoteMeta(qValue)
	test = strings.Replace(test, "\\*", ".*", -1)
	regex, err := regexp.Compile(test + "$")
	if err != nil {
		return false
	}
	switch op {
	case opNE:
		{
			return !regex.MatchString(oValue)
		}
	}
	return regex.MatchString(oValue)
}

func numericMatch(op string, qValue float64, oValue float64) bool {
	switch op {
	case opNE:
		{
			return qValue != oValue
		}
	case opGT:
		{
			return oValue > qValue
		}
	case opLT:
		{
			return oValue < qValue
		}
	}
	return qValue == oValue
}

func timeMatch(op string, qValue time.Time, oValue time.Time) bool {
	switch op {
	case opNE:
		{
			return !qValue.Equal(oValue)
		}
	case opGT:
		{
			return qValue.After(oValue)
		}
	case opLT:
		{
			return qValue.Before(oValue)
		}
	}
	return qValue.Equal(oValue)
}
