package utils

import "strings"

func StrContains(str1 string, str2 string) bool {
	return strings.Contains(strings.ToLower(str1), strings.ToLower(str2))
}
