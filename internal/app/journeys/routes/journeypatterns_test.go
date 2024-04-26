//go:build journeys_journeypattern_tests || journeys_tests || all_tests

package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestJourneyPatternMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.JourneyPattern
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.JourneyPattern{Route: nil}, map[string]string{"lineId": "1"}, false},
		{&model.JourneyPattern{Route: &model.Route{Line: nil}}, map[string]string{"lineId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"firstStopPointId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"lastStopPointId": "1"}, false},
		{&model.JourneyPattern{StopPoints: nil}, map[string]string{"stopPointId": "1"}, false},
	}

	for _, tc := range testCases {
		matches := journeyPatternMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestJourneyPatternsRoutes(t *testing.T) {
	testCases := []struct {
		target string
		items  []JourneyPattern
	}{
		{"/journey-patterns",
			[]JourneyPattern{
				getJourneyPatternMap()["047b0afc973ee2fd4fe92b128c3a932a"],
				getJourneyPatternMap()["65f51d2f85284af2fad1305c0ce71033"],
				getJourneyPatternMap()["9bc7403ad27267edbfbd63c3e92e5afa"],
				getJourneyPatternMap()["c01c71b0c9f456ba21f498a1dca54b3b"],
			},
		},
		{"/journey-patterns?lineId=1",
			[]JourneyPattern{
				getJourneyPatternMap()["047b0afc973ee2fd4fe92b128c3a932a"],
			},
		},
		{"/journey-patterns?firstStopPointId=4600&lastStopPointId=8149&stopPointId=8171",
			[]JourneyPattern{
				getJourneyPatternMap()["9bc7403ad27267edbfbd63c3e92e5afa"],
			},
		},
		{"/journey-patterns/047b0afc973ee2fd4fe92b128c3a932a?exclude-fields=name,stopPoints.municipality.url",
			[]JourneyPattern{
				{
					Url:             journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					LineUrl:         lineUrl("1"),
					RouteUrl:        routeUrl("1504270174600"),
					OriginStop:      stopPointUrl("7017"),
					DestinationStop: stopPointUrl("7015"),
					Direction:       "1",
					StopPoints: []StopPoint{
						{
							Url:        stopPointUrl("7017"),
							ShortName:  "7017",
							Name:       "Suupantori",
							Location:   "61.46546,23.64219",
							TariffZone: "B",
							Municipality: Municipality{
								ShortName: "604",
								Name:      "Pirkkala",
							},
						},
						{
							Url:        stopPointUrl("7015"),
							ShortName:  "7015",
							Name:       "Pirkkala",
							Location:   "61.4659,23.64734",
							TariffZone: "B",
							Municipality: Municipality{
								ShortName: "604",
								Name:      "Pirkkala",
							},
						},
					},
				},
			},
		},
		{"/journey-patterns/foobar",
			[]JourneyPattern{},
		},
		{"/journey-patterns?lastStopPointId=foobar",
			[]JourneyPattern{},
		},
		{"/journey-patterns?firstStopPointId=foobar",
			[]JourneyPattern{},
		},
		{"/journey-patterns?stopPointId=foobar",
			[]JourneyPattern{},
		},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectJourneyPatternRoutes(r, ctx)

		response := getJourneyPatternSuccessResponse(t, r, w, tc.target)

		dataSize := len(tc.items)
		if success := validateCommonResponseFields(t, response.Status, response.Data, uint16(dataSize)); !success {
			break
		}
		if len(response.Body) != dataSize {
			t.Errorf("expected %v, got %v", dataSize, len(response.Body))
			break
		}
		for i, l := range response.Body {
			if tc.items[i].Url != l.Url || tc.items[i].LineUrl != l.LineUrl ||
				tc.items[i].RouteUrl != l.RouteUrl || tc.items[i].OriginStop != l.OriginStop || tc.items[i].DestinationStop != l.DestinationStop ||
				tc.items[i].Direction != l.Direction {
				t.Errorf("expected %v, got %v", tc.items[i], l)
				break
			}

			for z := range tc.items[i].StopPoints {
				if tc.items[i].StopPoints[z] != l.StopPoints[z] {
					t.Errorf("expected %v, got %v", tc.items[i].StopPoints[z], l.StopPoints[z])
					break
				}
			}
		}
	}
}

type journeyPatternSuccessResponse struct {
	Status string           `json:"status"`
	Data   apiSuccessData   `json:"data"`
	Body   []JourneyPattern `json:"body"`
}

func getJourneyPatternSuccessResponse(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) journeyPatternSuccessResponse {
	serveHttp(t, r, w, target)

	var response journeyPatternSuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}
