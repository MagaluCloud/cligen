package strutils

import "strings"

// input: "/home/gfz/git/cligen/tmp-sdk/availabilityzones/az.go"
// output: "/tmp-sdk/availabilityzones/az.go"
func GetFileName(filePath string) string {
	spl := strings.Split(filePath, "/")
	for i, s := range spl {
		if s == "tmp-sdk" {
			return strings.Join(spl[i:], "/")
		}
	}
	return ""
}
