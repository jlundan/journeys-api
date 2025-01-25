//go:build journeys_stops_tests || journeys_tests || all_tests

package v1

import (
	"testing"
)

func TestGetExcludeFieldsQueryParameter(t *testing.T) {
	p := getExcludeFieldsQueryParameter(nil)

	if p != "" {
		t.Error("expected empty string")
	}
}
