package ggtfs

import "strings"

func StringIsNilOrEmpty(id *string) bool {
	return id == nil || strings.TrimSpace(*id) == ""
}
