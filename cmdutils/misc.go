package cmdutils

import (
	"strings"
)

// SplitBySpaces returns a string slice of words in s
func SplitBySpaces(s string) []string {

	s = strings.TrimSpace(s)
	sl := strings.Split(s, " ")
	return sl

}
