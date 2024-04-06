package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestMunicipalitiesMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Municipality
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
	}

	for _, tc := range testCases {
		matches := municipalityMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestMunicipalitiesRoutes(t *testing.T) {
	testCases := []struct {
		target string
		items  []Municipality
	}{
		{"/municipalities", []Municipality{
			{municipalityUrl("211"), "211", "Kangasala"},
			{municipalityUrl("604"), "604", "Pirkkala"},
			{municipalityUrl("837"), "837", "Tampere"},
			{municipalityUrl("980"), "980", "Ylöjärvi"},
		}},
		{"/municipalities?name=Pirkkala", []Municipality{
			{municipalityUrl("604"), "604", "Pirkkala"},
		}},
		{"/municipalities?name=Tampere&shortName=837", []Municipality{
			{municipalityUrl("837"), "837", "Tampere"},
		}},
		{"/municipalities?shortName=837&exclude-fields=name,shortName", []Municipality{
			{municipalityUrl("837"), "", ""},
		}},
		{"/municipalities/837", []Municipality{
			{municipalityUrl("837"), "837", "Tampere"},
		}},
		{"/municipalities/foobar", []Municipality{}},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectMunicipalityRoutes(r, ctx)

		response := getMunicipalitySuccessResponse(t, r, w, tc.target)

		dataSize := len(tc.items)
		if success := validateCommonResponseFields(t, response.Status, response.Data, uint16(dataSize)); !success {
			break
		}
		if len(response.Body) != dataSize {
			t.Errorf("expected %v, got %v", dataSize, len(response.Body))
			break
		}
		for i, l := range response.Body {
			if tc.items[i] != l {
				t.Errorf("expected %v, got %v", tc.items[i], l)
				break
			}
		}
	}
}

type municipalitySuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Municipality `json:"body"`
}

func getMunicipalitySuccessResponse(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) municipalitySuccessResponse {
	serveHttp(t, r, w, target)

	var response municipalitySuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}
