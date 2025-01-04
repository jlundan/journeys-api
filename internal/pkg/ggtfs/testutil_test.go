//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func compareStructs(expected, actual interface{}) (bool, string) {
	if expected == nil || actual == nil {
		if expected == actual {
			return true, ""
		}
		return false, fmt.Sprintf("One of the structs is nil: expected = %v, actual = %v", expected, actual)
	}

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	// Check if both values are pointers and dereference them.
	if expectedVal.Kind() == reflect.Ptr && actualVal.Kind() == reflect.Ptr {
		expectedVal = expectedVal.Elem()
		actualVal = actualVal.Elem()
	}

	// Ensure both are of the same type.
	if expectedVal.Type() != actualVal.Type() {
		return false, fmt.Sprintf("Type mismatch: expected type = %v, actual type = %v", expectedVal.Type(), actualVal.Type())
	}

	// Both values should be structs.
	if expectedVal.Kind() != reflect.Struct || actualVal.Kind() != reflect.Struct {
		return false, fmt.Sprintf("Both values must be structs: expected kind = %v, actual kind = %v", expectedVal.Kind(), actualVal.Kind())
	}

	var differences []string

	// Iterate through the fields of the struct.
	for i := 0; i < expectedVal.NumField(); i++ {
		fieldName := expectedVal.Type().Field(i).Name
		expectedField := expectedVal.Field(i)
		actualField := actualVal.Field(i)

		if !reflect.DeepEqual(expectedField.Interface(), actualField.Interface()) {
			differences = append(differences, fmt.Sprintf(
				"Field '%s': expected = %v, actual = %v",
				fieldName, expectedField.Interface(), actualField.Interface(),
			))
		}
	}

	// If no differences are found, the structs are equal.
	if len(differences) == 0 {
		return true, ""
	}

	// Return the differences as a formatted string.
	return false, fmt.Sprintf("Differences:\n%s", formatDifferences(differences))
}

// Helper function to format the differences as a readable string.
func formatDifferences(differences []string) string {
	return fmt.Sprintf("Found %d differences:\n- %s", len(differences), strings.Join(differences, "\n- "))
}

func tableToString(rows [][]string) string {
	var sb strings.Builder

	for _, row := range rows {
		sb.WriteString(strings.Join(row, ",") + "\n")
	}

	return sb.String()
}

func stringPtr(s string) *string {
	return &s
}

func handleEntityCreateResults[E GtfsEntity](t *testing.T, results []E, expectedResults []E) {
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d parsed structs, got %d", len(expectedResults), len(results))
		for i, result := range expectedResults {
			t.Logf("Expected result %d: %v", i, result)
		}

		for i, result := range results {
			t.Logf("Actual result %d: %v", i, result)
		}

		return
	}

	for i, expected := range expectedResults {
		isEqual, diff := compareStructs(expected, results[i])
		if !isEqual {
			t.Errorf("Struct comparison failed for entity %d:\n%s", i, diff)
		}
	}
}

func handleValidationResults(t *testing.T, results []Result, expectedResults []Result) {
	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d parsed structs, got %d", len(expectedResults), len(results))
		for i, result := range expectedResults {
			t.Logf("Expected result %d: (%s) %v", i, result.Code(), result)
		}

		for i, result := range results {
			t.Logf("Actual result %d: (%s) %v", i, result.Code(), result)
		}

		return
	}

	for i, expected := range expectedResults {
		isEqual, diff := compareStructs(expected, results[i])
		if !isEqual {
			t.Errorf("Struct comparison failed for entity %d:\n%s", i, diff)
		}
	}
}
