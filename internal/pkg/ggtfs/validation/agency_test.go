package validation

import (
	"testing"
)

func TestValidateAgency(t *testing.T) {
	tests := []struct {
		agency          Agency
		expectedResults []Result
	}{
		{
			agency: DummyCsvAgency{},
			expectedResults: []Result{
				MissingRequiredFieldResult{FieldName: "agency_name"},
				MissingRequiredFieldResult{FieldName: "agency_url"},
				MissingRequiredFieldResult{FieldName: "agency_timezone"},
			},
		},
	}

	for _, tt := range tests {
		results := ValidateAgency(tt.agency)
		for i, result := range results {
			if result != tt.expectedResults[i] {
				t.Errorf("Got %v, want %v", result, tt.expectedResults[i])
			}
		}
	}
}
