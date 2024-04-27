package ggtfs

import (
	"encoding/csv"
	"io"
	"strings"
)

func ReadHeaderRow(r *csv.Reader) (map[string]uint8, error) {
	row, err := r.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var headers = map[string]uint8{}
	for index, item := range row {
		headers[strings.TrimSpace(item)] = uint8(index)
	}
	return headers, nil
}

func ReadDataRow(r *csv.Reader) ([]string, error) {
	row, err := r.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var trimmedRow []string
	for _, item := range row {
		trimmedRow = append(trimmedRow, strings.TrimSpace(item))
	}
	return trimmedRow, nil
}

func StringArrayContainsItem(s []string, e string) bool {
	for _, a := range s {
		if a == e {
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
