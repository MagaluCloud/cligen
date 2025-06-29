package strutils

import "strings"

func FirstLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
