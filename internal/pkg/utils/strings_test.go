//go:build utils_tests || all_tests

package utils

import "testing"

func TestStrContains(t *testing.T) {
	str := "foobar"

	if !StrContains(str, "bar") {
		t.Error("expected true, got false")
	}

	if StrContains(str, "baz") {
		t.Error("expected false, got true")
	}
}
