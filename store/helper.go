package store

import (
	"regexp"
)

var stripAlphaNumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func sanitizeUsername(username string) string {
	return stripAlphaNumericRegex.ReplaceAllString(username, "")
}
