//go:build ggtfs_tests || all_tests

package ggtfs

import (
	"fmt"
	"sort"
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
