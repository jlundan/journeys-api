//go:build journeys_stops_tests || journeys_tests || all_tests

package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestStopPointMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.StopPoint
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.StopPoint{Municipality: nil}, map[string]string{"municipalityName": "foo"}, false},
		{&model.StopPoint{Municipality: nil}, map[string]string{"municipalityShortName": "foo"}, false},
	}

	for _, tc := range testCases {
		matches, err := stopPointMatchesConditions(tc.item, tc.conditions)
		if err != nil {
			t.Error(err)
		}
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestStopPointRoutes(t *testing.T) {
	var stopPointMap = getStopPointMap()
	testCases := []struct {
		target        string
		items         []StopPoint
		errorExpected bool
	}{
		{"/stop-points",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"], stopPointMap["4600"],
				stopPointMap["7015"], stopPointMap["7017"], stopPointMap["8149"], stopPointMap["8171"]}, false},
		{"/stop-points?tariffZone=B",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"], stopPointMap["4600"],
				stopPointMap["7015"], stopPointMap["7017"], stopPointMap["8171"]}, false,
		},
		{"/stop-points?municipalityName=Tampere",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"]}, false,
		},
		{"/stop-points?municipalityShortName=837",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"]}, false,
		},
		{"/stop-points?location=61.4445,23.87235",
			[]StopPoint{stopPointMap["3615"]}, false,
		},
		{"/stop-points?location=61,23:62,23.8",
			[]StopPoint{stopPointMap["7015"], stopPointMap["7017"]}, false,
		},
		{"/stop-points?name=Pirkkala",
			[]StopPoint{stopPointMap["7015"]}, false,
		},
		{"/stop-points?name=Sudenkorennontie&shortName=81",
			[]StopPoint{stopPointMap["8149"]}, false,
		},
		{"/stop-points?shortName=4600&exclude-fields=name,shortName,municipality.shortName", []StopPoint{
			{
				stopPointUrl("4600"),
				"",
				"",
				"61.47561,23.97756",
				"B",
				Municipality{
					Url:  municipalityUrl("211"),
					Name: "Kangasala",
				},
			},
		}, false},
		{"/stop-points/4600", []StopPoint{
			stopPointMap["4600"],
		}, false},
		{"/stop-points/foobar", []StopPoint{}, false},
		{"/stop-points?location=6123:62,23.8", []StopPoint{}, true},
		{"/stop-points?location=61,23:6223.8", []StopPoint{}, true},
		{"/stop-points?location=A1,23:62,23.8", []StopPoint{}, true},
		{"/stop-points?location=61,23:62,23.X", []StopPoint{}, true},
		{"/stop-points?location=61,2X:62,23.8", []StopPoint{}, true},
		{"/stop-points?location=61,23:6X,23.8", []StopPoint{}, true},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectStopPointRoutes(r, ctx)

		if tc.errorExpected {
			_, err := getStopPointErrorResponse(t, r, w, tc.target)
			if err != nil {
				t.Error(err)
			}
			continue
		}

		gotResponse := getStopPointSuccessResponse(t, r, w, tc.target)
		expectedResponse := stopPointSuccessResponse{
			Status: "success",
			Data: apiSuccessData{
				Headers: apiHeaders{
					Paging: apiHeadersPaging{
						StartIndex: 0,
						PageSize:   uint16(len(tc.items)),
						MoreData:   false,
					},
				},
			},
			Body: tc.items,
		}

		var diffs []FieldDiff
		initialTag := fmt.Sprintf("%v:Response", tc.target)
		var err = compareVariables(expectedResponse, gotResponse, initialTag, &diffs, false)
		if err != nil {
			t.Error(err)
			break
		}

		if len(diffs) > 0 {
			printFieldDiffs(t, diffs)
			break
		}
	}
}

func getStopPointMap() map[string]StopPoint {
	result := make(map[string]StopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality Municipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getMunicipalityMap()["837"]},
	}

	for _, tc := range stopPoints {
		result[tc.id] = StopPoint{
			Url:          stopPointUrl(tc.id),
			ShortName:    tc.id,
			Name:         tc.name,
			Location:     tc.location,
			TariffZone:   tc.tariffZone,
			Municipality: tc.municipality,
		}
	}

	return result
}

type stopPointSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []StopPoint    `json:"body"`
}

func getStopPointSuccessResponse(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) stopPointSuccessResponse {
	serveHttp(t, r, w, target)

	var response stopPointSuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}

func getStopPointErrorResponse(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) (apiFailResponse, error) {
	serveHttp(t, r, w, target)

	var response apiFailResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return apiFailResponse{}, err
	}

	return response, nil
}
