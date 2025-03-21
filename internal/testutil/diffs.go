package testutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type FieldDiff struct {
	Tag      string
	Expected interface{}
	Got      interface{}
}

func CompareVariablesAndPrintResults(t *testing.T, expected, got interface{}, tag string) {
	var diffs []FieldDiff
	err := CompareVariables(expected, got, tag, &diffs, false)
	if err != nil {
		t.Error(err)
	}
	if len(diffs) > 0 {
		PrintFieldDiffs(t, diffs)
	}
}

func CompareVariables(expected, got interface{}, tag string, diffs *[]FieldDiff, verbose bool) error {
	expErr, expIsErr := expected.(error)
	gotErr, gotIsErr := got.(error)

	if expIsErr || gotIsErr {
		// If one is an error and the other isn't, they are different.
		if expIsErr != gotIsErr {
			*diffs = append(*diffs, FieldDiff{
				Tag:      tag,
				Expected: expected,
				Got:      got,
			})
			return nil
		}

		// Use errors.Is() to compare errors correctly.
		if !errors.Is(gotErr, expErr) {
			*diffs = append(*diffs, FieldDiff{
				Tag:      tag,
				Expected: expErr.Error(),
				Got:      gotErr.Error(),
			})
		}
		return nil
	}

	jsonExp, err1 := json.Marshal(expected)
	if err1 != nil {
		return err1
	}
	jsonGot, err2 := json.Marshal(got)
	if err2 != nil {
		return err2
	}

	if reflect.TypeOf(expected) != reflect.TypeOf(got) {
		return errors.New("structs are of different types")
	}

	va := reflect.ValueOf(expected)
	vb := reflect.ValueOf(got)

	if verbose {
		fmt.Println(fmt.Sprintf("---%v [%v]---\n--> a: %v\n--> b: %v", tag, va.Kind(), string(jsonExp), string(jsonGot)))
	}

	if va.Kind() == reflect.Ptr && vb.Kind() == reflect.Ptr {
		if va.IsNil() || vb.IsNil() {
			if va.IsNil() != vb.IsNil() {
				*diffs = append(*diffs, FieldDiff{
					Tag:      tag,
					Expected: va.Interface(),
					Got:      vb.Interface(),
				})
			}
			return nil
		}
		return CompareVariables(va.Elem().Interface(), vb.Elem().Interface(), tag, diffs, verbose)
	}

	if va.Kind() == reflect.Slice && vb.Kind() == reflect.Slice {
		if va.Len() != vb.Len() {
			*diffs = append(*diffs, FieldDiff{
				Tag:      tag,
				Expected: va.Interface(),
				Got:      vb.Interface(),
			})
			return nil
		}

		for x := 0; x < va.Len(); x++ {
			err := CompareVariables(va.Index(x).Interface(), vb.Index(x).Interface(), fmt.Sprintf("%v.[%v]", tag, x), diffs, verbose)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if va.Kind() == reflect.Struct && vb.Kind() == reflect.Struct {
		for i := 0; i < va.NumField(); i++ {
			fieldA := va.Field(i)
			fieldB := vb.Field(i)
			err := CompareVariables(fieldA.Interface(), fieldB.Interface(), fmt.Sprintf("%v.%v", tag, va.Type().Field(i).Name), diffs, verbose)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if !reflect.DeepEqual(expected, got) {
		*diffs = append(*diffs, FieldDiff{
			Tag:      tag,
			Expected: expected,
			Got:      got,
		})
		return nil
	}

	if verbose {
		fmt.Println("Compare OK")
	}

	return nil
}

func PrintFieldDiffs(t *testing.T, diffs []FieldDiff) {
	for _, diff := range diffs {
		var expectedJSON, expectedErr = json.Marshal(diff.Expected)
		if expectedErr != nil {
			t.Fatalf("Failed to marshal expected JSON: %v", expectedErr)
		}
		var gotJSON, gotErr = json.Marshal(diff.Got)
		if gotErr != nil {
			t.Fatalf("Failed to marshal expected JSON: %v", gotErr)
		}

		t.Error(fmt.Sprintf("Case: %v\n expected: \n%v \ngot: \n%v", diff.Tag, string(expectedJSON), string(gotJSON)))
	}
}
