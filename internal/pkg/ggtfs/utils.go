package ggtfs

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func ReadHeaderRow(r *csv.Reader, validHeaders []string) (map[string]int, error) {
	row, err := r.Read()
	if err != nil {
		return nil, err
	}

	validHeaderSet := make(map[string]struct{}, len(validHeaders))
	for _, header := range validHeaders {
		validHeaderSet[header] = struct{}{}
	}

	headers := make(map[string]int)
	encounteredHeaders := make(map[string]bool)

	for index, headerName := range row {
		headerName = strings.TrimSpace(headerName)

		if _, isValid := validHeaderSet[headerName]; !isValid {
			continue
		}

		if encounteredHeaders[headerName] {
			return nil, fmt.Errorf("duplicate header found: %s", headerName)
		}

		encounteredHeaders[headerName] = true
		headers[headerName] = index
	}
	return headers, nil
}

func StringArrayContainsItem(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

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
