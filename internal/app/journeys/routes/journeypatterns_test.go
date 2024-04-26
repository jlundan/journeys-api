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
					Journeys: []JourneyPatternJourney{
						{
							Url:               "http://localhost:5678/journeys/7020205685",
							JourneyPatternUrl: "http://localhost:5678/journey-patterns/047b0afc973ee2fd4fe92b128c3a932a",
							DepartureTime:     "14:43:00",
							ArrivalTime:       "14:44:45",
							HeadSign:          "Vatiala",
							DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
							DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
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

		gotResponse := getJourneyPatternSuccessResponse(t, r, w, tc.target)
		expectedResponse := journeyPatternSuccessResponse{
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
