//go:build utils_tests || all_tests

package utils

import (
	"errors"
	"testing"
)

var fakeMarshal, originalMarshal func(_ interface{}) ([]byte, error)
var fakeUnmarshal, originalUnmarshal func(data []byte, v interface{}) error

func init() {
	fakeMarshal = func(_ interface{}) ([]byte, error) {
		return []byte{}, errors.New("marshalling failed")
	}

	fakeUnmarshal = func(data []byte, v interface{}) error {
		return errors.New("marshalling failed")
	}
}

func TestToInterfaceMapViaJSON(t *testing.T) {
	o := testObj{
		Field1: "foo",
		Sub:    subObj{SubField1: "bar"},
	}

	m, err := ToInterfaceMapViaJSON(o)
	if err != nil {
		t.Error(err)
	}

	if m["Field1"] != "foo" || m["Sub"] == nil {
		t.Error("expected field1 to be foo and sub not to be nil")
	}

	if m["Sub"].(map[string]interface{})["SubField1"] != "bar" {
		t.Error("expected SubField1 to be bar")
	}

	originalMarshal = jsonMarshal
	jsonMarshal = fakeMarshal
	_, err = ToInterfaceMapViaJSON(o)
	if err == nil {
		t.Error("expected to receive an error")
	}
	jsonMarshal = originalMarshal

	originalUnmarshal = jsonUnmarshal
	jsonUnmarshal = fakeUnmarshal
	_, err = ToInterfaceMapViaJSON(o)
	if err == nil {
		t.Error("expected to receive an error")
	}
	jsonUnmarshal = originalUnmarshal
}

func TestFilterObject(t *testing.T) {
	o := testObj{
		Field1: "foo",
		Sub:    subObj{SubField1: "bar"},
	}

	filtered, err := FilterObject(o, "Field1")
	if err != nil {
		t.Error(err)
	}

	if _, ok := filtered.(map[string]interface{})["Field1"]; ok {
		t.Error("expected Field1 to be gone")
	}

	filtered, err = FilterObject(o, "")
	if err != nil {
		t.Error(err)
	}

	originalMarshal = jsonMarshal
	jsonMarshal = fakeMarshal
	_, err = FilterObject(o, "")
	if err == nil {
		t.Error("expected to receive an error")
	}
	jsonMarshal = originalMarshal
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

	m, err := ToInterfaceMapViaJSON(o)
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
