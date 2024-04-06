package ggtfs

import (
	"fmt"
	"testing"
)

func TestTableToString(t *testing.T) {
	testCases := []struct {
		rows     [][]string
		expected string
	}{
		{
			rows: [][]string{
				{"agency_name"},
				{","},
			},
			expected: "agency_name\n,\n",
		},
		{
			rows: [][]string{
				{"agency_name"},
				{" "},
			},
			expected: "agency_name\n \n",
		},
	}
	for _, tc := range testCases {
		if tableToString(tc.rows) != tc.expected {
			t.Error(fmt.Sprintf("expected %s, got %s", tc.expected, tableToString(tc.rows)))
		}
	}
}
