//go:build utils_tests || journeys_tests || all_tests

package v1

import (
	"testing"
)

func TestRemoveExcludedFields(t *testing.T) {
	m := removeExcludedFields([]map[string]interface{}{}, "")

	if len(m) != 0 {
		t.Error("expected empty map")
	}
}

func TestConvertToStringAnyMap(t *testing.T) {
	type testStruct struct {
		Field1 string
		Field2 bool
	}

	tests := []struct {
		name    string
		input   any
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Valid struct",
			input: testStruct{
				Field1: "value1",
				Field2: true,
			},
			want: map[string]any{
				"Field1": "value1",
				"Field2": true,
			},
			wantErr: false,
		},
		{
			name:    "Invalid input",
			input:   make(chan int),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToStringAnyMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToStringAnyMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalMaps(got, tt.want) {
				t.Errorf("convertToStringAnyMap() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test that non-structs cannot be serialized
	// This is sort of unnecessary since convertToStringAnyMap is only used with APIEntities which are structs,
	// but the test is here for completeness
	_, err := convertToStringAnyMap("foo-bar")
	if err == nil {
		t.Error("expected error")
	}
}

func equalMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func TestDeleteProperty(t *testing.T) {
	o := testObj{
		Field1: "foo",
		Sub: subObj{
			SubField1: "bar",
			SubField2: []subSubObj{
				{SubSubField1: "baz"},
			},
		},
	}

	m, err := convertToStringAnyMap(o)
	if err != nil {
		t.Error(err)
	}

	deleteProperty(m, "Sub.SubField1")
	if _, ok := m["Sub"].(map[string]interface{})["SubField1"]; ok {
		t.Error("expected SubField1 to be gone")
	}

	deleteProperty(m, "Sub.SubField2.SubSubField1")
	sub := m["Sub"].(map[string]interface{})
	subField2 := sub["SubField2"].([]interface{})[0]
	if _, ok := subField2.(map[string]interface{})["SubSubField1"]; ok {
		t.Error("expected SubSubField1 to be gone")
	}
}

type subSubObj struct {
	SubSubField1 string
}

type subObj struct {
	SubField1 string
	SubField2 []subSubObj
}

type testObj struct {
	Field1 string
	Sub    subObj
}
