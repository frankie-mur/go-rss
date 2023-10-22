package validator

import (
	"regexp"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var RssUrlRX = regexp.MustCompile("https?://[^/]+/(?:rss|feed|atom)(?:/.*\\.xml)?|https?://[^/]+/.*\\.xml\n")

// Define a new Validator type which contains a map of validation errors for our
// form fields.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() returns true if the FieldErrors and NonFieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Matches() returns true if a value matches a provided compiled regular
// expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
