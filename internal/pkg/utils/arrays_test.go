//go:build utils_tests || all_tests

package utils

import "testing"

func TestStringArrayContainsItem(t *testing.T) {
	arr := []string{"item1", "item2"}

	if !StringArrayContainsItem(arr, "item1") {
		t.Error("expected true, got false")
	}

	if StringArrayContainsItem(arr, "foo") {
		t.Error("expected false, got true")
	}
}
