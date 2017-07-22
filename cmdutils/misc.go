package cmdutils

import (
	"strings"
)

func SplitBySpaces(s string) []string {

	s = strings.TrimSpace(s)
	sl := strings.Split(s, " ")
	return sl

}
