package util

import "strings"

// Core holds various utility methods
type Core struct{}

// StrReplace replaces one string with another.
// This is methodized so it can be used in a template.
func (s Core) StrReplace(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}
