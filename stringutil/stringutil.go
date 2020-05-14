package stringutil

import "strings"

func EqualIgnoreCase(s1, s2 string) bool {
	return strings.ToUpper(s1) == strings.ToUpper(s2)
}
