package ggtfs

import (
	"strings"
)

func tableToString(rows [][]string) string {
	var sb strings.Builder

	for _, row := range rows {
		sb.WriteString(strings.Join(row, ",") + "\n")
	}

	return sb.String()
}

func toSet[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}
