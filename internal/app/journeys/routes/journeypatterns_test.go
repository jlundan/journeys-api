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

func getJourneyPatternMap() map[string]JourneyPattern {
	result := make(map[string]JourneyPattern)

	journeyPatterns := []struct {
		id              string
		line            string
		route           string
		originStop      string
		destinationStop string
		name            string
		direction       string
		stopPoints      []StopPoint
		journeys        []JourneyPatternJourney
	}{
		{"047b0afc973ee2fd4fe92b128c3a932a", "1", "1504270174600", "7017",
			"7015", "Suupantori - Pirkkala", "1",
			[]StopPoint{getStopPointMap()["7017"], getStopPointMap()["7015"]},
			[]JourneyPatternJourney{
				{
					Url:               "http://localhost:5678/journeys/7020205685",
					JourneyPatternUrl: "http://localhost:5678/journey-patterns/047b0afc973ee2fd4fe92b128c3a932a",
					DepartureTime:     "14:43:00",
					ArrivalTime:       "14:44:45",
					HeadSign:          "Vatiala",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
			}},

		{"65f51d2f85284af2fad1305c0ce71033", "3A", "1517136151028", "3615",
			"3607", "Näyttelijänkatu - Lavastajanpolku", "0",
			[]StopPoint{getStopPointMap()["3615"], getStopPointMap()["3607"]},
			[]JourneyPatternJourney{
				{
					Url:               "http://localhost:5678/journeys/7024545685",
					JourneyPatternUrl: "http://localhost:5678/journey-patterns/65f51d2f85284af2fad1305c0ce71033",
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					HeadSign:          "Lentävänniemi",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
			}},

		{"9bc7403ad27267edbfbd63c3e92e5afa", "1A", "1501146007035", "4600",
			"8149", "Vatiala - Sudenkorennontie", "0",
			[]StopPoint{getStopPointMap()["4600"], getStopPointMap()["8171"], getStopPointMap()["8149"]},
			[]JourneyPatternJourney{
				{
					Url:               "http://localhost:5678/journeys/7020295685",
					JourneyPatternUrl: "http://localhost:5678/journey-patterns/9bc7403ad27267edbfbd63c3e92e5afa",
					DepartureTime:     "06:30:00",
					ArrivalTime:       "06:32:30",
					HeadSign:          "Lentoasema",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
			}},

		{"c01c71b0c9f456ba21f498a1dca54b3b", "-1", "111111111", "3615",
			"7017", "Näyttelijänkatu - Suupantori", "0",
			[]StopPoint{getStopPointMap()["3615"], getStopPointMap()["7017"]},
			[]JourneyPatternJourney{
				{
					Url:               "http://localhost:5678/journeys/111111111",
					JourneyPatternUrl: "http://localhost:5678/journey-patterns/c01c71b0c9f456ba21f498a1dca54b3b",
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					HeadSign:          "Foobar",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{},
				},
			}},
	}

	for _, tc := range journeyPatterns {
		result[tc.id] = JourneyPattern{
			Url:             journeyPatternUrl(tc.id),
			LineUrl:         lineUrl(tc.line),
			RouteUrl:        routeUrl(tc.route),
			OriginStop:      stopPointUrl(tc.originStop),
			DestinationStop: stopPointUrl(tc.destinationStop),
			Direction:       tc.direction,
			Name:            tc.name,
			StopPoints:      tc.stopPoints,
			Journeys:        tc.journeys,
		}
	}

	return result
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
