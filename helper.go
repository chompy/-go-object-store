package main

import (
	"regexp"
	"strings"
)

var stripAlphaNumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func stringMatch(qValue string, oValue string) bool {
	test := regexp.QuoteMeta(qValue)
	test = strings.Replace(test, "\\*", ".*", -1)
	regex, err := regexp.Compile(test + "$")
	if err != nil {
		return false
	}
	return regex.MatchString(oValue)
}

func sanitizeUsername(username string) string {
	return stripAlphaNumericRegex.ReplaceAllString(username, "")
}
