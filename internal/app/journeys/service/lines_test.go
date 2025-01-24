package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"testing"
)

func TestLineMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Line
		conditions map[string]string
		expected   bool
	}{
		{
			nil,
			nil,
			false,
		},
		{
			&model.Line{Name: "1", Description: "Vatiala - Pirkkala"},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		matches := lineMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}
