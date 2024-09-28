//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func checkErrors(expected []string, actual []error, t *testing.T) {
	if len(expected) == 0 && len(actual) != 0 {
		t.Error(fmt.Sprintf("expected zero errors, got %v", len(actual)))
		for _, a := range actual {
			fmt.Println(a)
		}
		return
	}

	if len(expected) != 0 && len(actual) == 0 {
		t.Error(fmt.Sprintf("expected %v errors, got zero", len(expected)))
		for _, a := range actual {
			fmt.Println(a)
		}
		return
	}

	if len(expected) != len(actual) {
		t.Error(fmt.Sprintf("expected %v errors, got %v", len(expected), len(actual)))
		for _, a := range actual {
			fmt.Println(a)
		}
		return
	}

	sort.Slice(actual, func(x, y int) bool {
		return actual[x].Error() < actual[y].Error()
	})

	sort.Slice(expected, func(x, y int) bool {
		return expected[x] < expected[y]
	})

	for i, a := range actual {
		if a.Error() != expected[i] {
			t.Error(fmt.Sprintf("expected error %s, got %s", expected[i], a.Error()))
		}
	}
}

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

type ggtfsTestCase struct {
	csvRows         [][]string
	expectedStructs []interface{}
	expectedErrors  []string
}

type LoadFunction func(reader *csv.Reader) ([]interface{}, []error)
type ValidateFunction func(entities []interface{}) []error

type ParseResult struct {
	Entities []interface{}
	Errors   []error
}

// LoadAndValidateGTFS is a generic function to load and validate GTFS entities while allowing partial success.
func loadAndValidateGTFS(csvReader *csv.Reader, loadFunc LoadFunction, validateFunc ValidateFunction, strictMode bool) ParseResult {
	// Load the entities using the provided load function
	entities, parseErrors := loadFunc(csvReader)
	validationErrors := validateFunc(entities)

	// If strict mode is enabled, combine parse errors and validation errors and return an empty set of entities
	if strictMode && (len(parseErrors) > 0 || len(validationErrors) > 0) {
		return ParseResult{
			Entities: nil,
			Errors:   append(parseErrors, validationErrors...),
		}
	}

	// Otherwise, return the parsed entities and the errors separately
	return ParseResult{
		Entities: entities,
		Errors:   append(parseErrors, validationErrors...),
	}
}

func runGenericGTFSParseTest(t *testing.T, testName string, loadFunc LoadFunction, validateFunc ValidateFunction, strictMode bool, testCases map[string]ggtfsTestCase) {
	for tcName, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", testName, tcName), func(t *testing.T) {
			result := loadAndValidateGTFS(csv.NewReader(strings.NewReader(tableToString(tc.csvRows))), loadFunc, validateFunc, strictMode)

			// Sort errors for consistent comparison
			sort.Slice(result.Errors, func(x, y int) bool {
				return result.Errors[x].Error() < result.Errors[y].Error()
			})
			sort.Strings(tc.expectedErrors)

			// Check error count and contents
			if len(result.Errors) != len(tc.expectedErrors) {
				t.Errorf("Expected %d errors, got %d", len(tc.expectedErrors), len(result.Errors))
				for _, e := range result.Errors {
					t.Logf("Actual error: %s", e.Error())
				}
				t.Logf("Expected errors: %v", tc.expectedErrors)
				return
			}

			for i, e := range result.Errors {
				if e.Error() != tc.expectedErrors[i] {
					t.Errorf("Expected error %q, got %q", tc.expectedErrors[i], e.Error())
				}
			}

			// Check that the parsed entities match the expected entities if they are provided
			if len(tc.expectedStructs) > 0 {
				if len(result.Entities) != len(tc.expectedStructs) {
					t.Errorf("Expected %d parsed structs, got %d", len(tc.expectedStructs), len(result.Entities))
					return
				}

				for i, expected := range tc.expectedStructs {
					// Use the provided compareStructs function to compare the actual struct with the expected one.
					isEqual, diff := compareStructs(expected, result.Entities[i])
					if !isEqual {
						t.Errorf("Struct comparison failed for entity %d:\n%s", i, diff)
					}
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
