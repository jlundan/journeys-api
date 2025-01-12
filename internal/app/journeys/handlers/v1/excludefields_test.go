//go:build utils_tests || all_tests

package v1

import (
	"testing"
)

//func TestToInterfaceMapViaJSON(t *testing.T) {
//	o := testObj{
//		Field1: "foo",
//		Sub:    subObj{SubField1: "bar"},
//	}
//
//	m, err := convertToMap(o)
//	if err != nil {
//		t.Error(err)
//	}
//
//	if m["Field1"] != "foo" || m["Sub"] == nil {
//		t.Error("expected field1 to be foo and sub not to be nil")
//	}
//
//	if m["Sub"].(map[string]interface{})["SubField1"] != "bar" {
//		t.Error("expected SubField1 to be bar")
//	}
//
//	_, err = convertToMap(o)
//	if err == nil {
//		t.Error("expected to receive an error")
//	}
//
//	_, err = convertToMap(o)
//	if err == nil {
//		t.Error("expected to receive an error")
//	}
//}

func TestFilterObject(t *testing.T) {
	o := testObj{
		Field1: "foo",
		Sub:    subObj{SubField1: "bar"},
	}

	m, err := convertToMap(o)
	filtered, err := filterMap(m, "Field1")
	if err != nil {
		t.Error(err)
	}

	if _, ok := filtered.(map[string]interface{})["Field1"]; ok {
		t.Error("expected Field1 to be gone")
	}

	m2, err := convertToMap(o)
	filtered, err = filterMap(m2, "")
	if err != nil {
		t.Error(err)
	}

	//type Test struct {
	//	Ch chan int
	//}
	//test := Test{Ch: make(chan int)}
	//_, err = filterMap(test, "")
	//if err == nil {
	//	t.Error("expected to receive an error")
	//}
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

	m, err := convertToMap(o)
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
