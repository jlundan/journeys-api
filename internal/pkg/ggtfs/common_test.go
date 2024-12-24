//go:build ggtfs_tests || all_tests || ggtfs_tests_common

package ggtfs

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestGetRowValueForHeaderName(t *testing.T) {
	headerIndex := map[string]int{
		"test1": 0,
		"test2": 1,
		"test3": 2,
	}
	if *getRowValueForHeaderName([]string{"foo", "bar", "baz"}, headerIndex, "test1") != "foo" {
		t.Error("expected test1")
	}

	if *getRowValueForHeaderName([]string{"foo", "bar", "baz"}, headerIndex, "test2") != "bar" {
		t.Error("expected bar")
	}

	if *getRowValueForHeaderName([]string{"foo", "bar", "baz"}, headerIndex, "test3") != "baz" {
		t.Error("expected baz")
	}

	if getRowValueForHeaderName([]string{"foo", "bar", "baz"}, headerIndex, "test4") != nil {
		t.Error("expected nil")
	}
}

func TestGetHeaderIndex(t *testing.T) {
	validHeaders := []string{"test1", "test2", "test3"}
	headers := [][]string{
		{"test1", "test2", "test3", "test4"},
	}
	headerIndex, errors := getHeaderIndex(csv.NewReader(strings.NewReader(tableToString(headers))), validHeaders)

	if len(errors) != 0 {
		t.Error("expected zero errors")
	}

	if headerIndex["test1"] != 0 || headerIndex["test2"] != 1 || headerIndex["test3"] != 2 || headerIndex["test4"] != -1 {
		t.Error("expected four items in the index")
	}

	headers = [][]string{
		{"test1", "test2", "test2"},
	}
	headerIndex, errors = getHeaderIndex(csv.NewReader(strings.NewReader(tableToString(headers))), validHeaders)

	if len(errors) != 1 && errors[0].Error() != "duplicate header found: test2" {
		t.Error("expected duplicate header error")
	}

	if len(headerIndex) != 2 {
		t.Error("expected two items in the index")
	}

	headers = [][]string{{}}
	headerIndex, errors = getHeaderIndex(csv.NewReader(strings.NewReader(tableToString(headers))), validHeaders)

	if len(errors) != 0 {
		t.Error("expected zero errors")
	}

	if len(headerIndex) != 0 {
		t.Error("expected zero items in the index")
	}

	headerIndex, errors = getHeaderIndex(csv.NewReader(strings.NewReader("\"field1,field2\n")), validHeaders)
	if len(errors) != 1 {
		t.Error("expected one errors")
	}

}

func TestLoadEntitiesFromCSV(t *testing.T) {
	validHeaders := []string{"test1", "test2", "test3"}
	headers := [][]string{
		{"test1", "test2", "test3", "test3"},
		{" "},
		{","},
		{"", ""},
		{" ", " "},
	}
	_, errors := LoadEntitiesFromCSV(csv.NewReader(strings.NewReader(tableToString(headers))), validHeaders, dummyEntityCreator, "test.csv")

	if len(errors) != 5 {
		t.Error("expected four errors")
	}

	if errors[0].Error() != "test.csv: duplicate header name: test3" {
		t.Error("unexpected error")
	}
	if errors[1].Error() != "test.csv: record on line 2: wrong number of fields" {
		t.Error("unexpected error")
	}
	if errors[2].Error() != "test.csv: record on line 3: wrong number of fields" {
		t.Error("unexpected error")
	}
	if errors[3].Error() != "test.csv: record on line 4: wrong number of fields" {
		t.Error("unexpected error")
	}
	if errors[4].Error() != "test.csv: record on line 5: wrong number of fields" {
		t.Error("unexpected error")
	}

	headers = [][]string{}
	index, errors := LoadEntitiesFromCSV(csv.NewReader(strings.NewReader(tableToString(headers))), validHeaders, dummyEntityCreator, "test.csv")

	if len(errors) != 0 {
		t.Error("expected zero errors")
	}

	if len(index) != 0 {
		t.Error("expected zero items in the index")
	}
}

func dummyEntityCreator(row []string, headers map[string]int, lineNumber int) interface{} {
	return nil
}
