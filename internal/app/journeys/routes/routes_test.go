//go:build journeys_routes_tests || journeys_tests || all_tests

package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestRoutesMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Route
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.Route{Line: nil}, map[string]string{"lineId": "1"}, false},
	}

	for _, tc := range testCases {
		matches := routeMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestGetRoutes(t *testing.T) {
	rm := getRouteMap()

	testCases := []struct {
		target string
		items  []Route
	}{
		{"/routes",
			[]Route{rm[routeUrl("111111111")], rm[routeUrl("1501146007035")], rm[routeUrl("1504270174600")], rm[routeUrl("1517136151028")]},
		},
		{"/routes/1501146007035",
			[]Route{rm[routeUrl("1501146007035")]},
		},
		{"/routes?lineId=1A",
			[]Route{rm[routeUrl("1501146007035")]},
		},
		{"/routes/foobar",
			[]Route{},
		},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectRouteRoutes(r, ctx)

		gotResponse := sendRoutesRequest(t, r, w, tc.target)
		expectedResponse := routeSuccessResponse{
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

type routeSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Route        `json:"body"`
}

func sendRoutesRequest(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) routeSuccessResponse {
	serveHttp(t, r, w, target)

	var response routeSuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}

func getRouteMap() map[string]Route {
	result := make(map[string]Route)

	routes := []Route{
		{
			Url:        routeUrl("111111111"),
			LineUrl:    lineUrl("-1"),
			Name:       "Näyttelijänkatu - Suupantori",
			Projection: "6144444,2387262:17,-14",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             journeyPatternUrl("c01c71b0c9f456ba21f498a1dca54b3b"),
					OriginStop:      stopPointUrl("3615"),
					DestinationStop: stopPointUrl("7017"),
					Name:            "Näyttelijänkatu - Suupantori",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               journeyUrl("111111111"),
					JourneyPatternUrl: journeyPatternUrl("c01c71b0c9f456ba21f498a1dca54b3b"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{},
				},
			},
		},
		{
			Url:        routeUrl("1501146007035"),
			LineUrl:    lineUrl("1A"),
			Name:       "Vatiala - Sudenkorennontie",
			Projection: "6147590,2397764:-104,32",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             journeyPatternUrl("9bc7403ad27267edbfbd63c3e92e5afa"),
					OriginStop:      stopPointUrl("4600"),
					DestinationStop: stopPointUrl("8149"),
					Name:            "Vatiala - Sudenkorennontie",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               journeyUrl("7020295685"),
					JourneyPatternUrl: journeyPatternUrl("9bc7403ad27267edbfbd63c3e92e5afa"),
					DepartureTime:     "06:30:00",
					ArrivalTime:       "06:32:30",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
		{
			Url:        routeUrl("1504270174600"),
			LineUrl:    lineUrl("1"),
			Name:       "Suupantori - Pirkkala",
			Projection: "6146557,2364305:-19,-148:-17,-191",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					OriginStop:      stopPointUrl("7017"),
					DestinationStop: stopPointUrl("7015"),
					Name:            "Suupantori - Pirkkala",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               journeyUrl("7020205685"),
					JourneyPatternUrl: journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					DepartureTime:     "14:43:00",
					ArrivalTime:       "14:44:45",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
		{
			Url:        routeUrl("1517136151028"),
			LineUrl:    lineUrl("3A"),
			Name:       "Näyttelijänkatu - Lavastajanpolku",
			Projection: "6144444,2387262:17,-14",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             journeyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					OriginStop:      stopPointUrl("3615"),
					DestinationStop: stopPointUrl("3607"),
					Name:            "Näyttelijänkatu - Lavastajanpolku",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               journeyUrl("7024545685"),
					JourneyPatternUrl: journeyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
	}

	for _, route := range routes {
		result[route.Url] = route
	}

	return result
}
