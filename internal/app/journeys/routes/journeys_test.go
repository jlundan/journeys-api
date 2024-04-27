//go:build journeys_journeys_tests || journeys_tests || all_tests

package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestJourneyMatchesConditions(t *testing.T) {
	emptyJourney := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: nil, DayTypes: nil,
		DayTypeExceptions: nil, Calls: nil, Line: nil, JourneyPattern: nil, ValidFrom: "19700101",
		ValidTo: "20300101", Route: nil, ArrivalTime: "", DepartureTime: "", ActivityId: "",
	}
	journeyWithEmptyCallArr := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: nil, DayTypes: nil,
		DayTypeExceptions: nil, Calls: make([]*model.JourneyCall, 0), Line: nil, JourneyPattern: nil, ValidFrom: "19700101",
		ValidTo: "20300101", Route: nil, ArrivalTime: "", DepartureTime: "", ActivityId: "",
	}
	journeyWithEmptyDtArr := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: nil,
		DayTypes: make([]string, 0), DayTypeExceptions: nil, Calls: nil, Line: nil, JourneyPattern: nil, ValidFrom: "19700101",
		ValidTo: "20300101", Route: nil, ArrivalTime: "", DepartureTime: "", ActivityId: "",
	}
	invalidJourneyLowerSide := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: nil,
		DayTypes: make([]string, 0), DayTypeExceptions: nil, Calls: nil, Line: nil, JourneyPattern: nil, ValidFrom: "20300101",
		ValidTo: "20300101", Route: nil, ArrivalTime: "", DepartureTime: "", ActivityId: "",
	}
	invalidJourneyUpperSide := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: nil,
		DayTypes: make([]string, 0), DayTypeExceptions: nil, Calls: nil, Line: nil, JourneyPattern: nil, ValidFrom: "20300101",
		ValidTo: "20300101", Route: nil, ArrivalTime: "", DepartureTime: "", ActivityId: "",
	}
	validJourney := model.Journey{Id: "", HeadSign: "", Direction: "", WheelchairAccessible: false, GtfsInfo: &model.JourneyGtfsInfo{TripId: "1111"},
		DayTypes: []string{"monday", "tuesday"}, DayTypeExceptions: nil, Calls: []*model.JourneyCall{
			{StopPoint: &model.StopPoint{ShortName: "1"}}, {StopPoint: &model.StopPoint{ShortName: "2"}}, {StopPoint: &model.StopPoint{ShortName: "3"}}},
		Line: &model.Line{Description: "Foobar", Name: "1A"}, JourneyPattern: &model.JourneyPattern{Id: "123"},
		ValidFrom: "19700101", ValidTo: "20300101", Route: &model.Route{Id: "123"}, ArrivalTime: "02:00", DepartureTime: "01:00", ActivityId: "",
	}

	testCases := []struct {
		item       *model.Journey
		conditions map[string]string
		expected   bool
	}{
		{item: &emptyJourney, conditions: map[string]string{"lineId": "1"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"routeId": "1"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"journeyPatternId": "1"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"dayTypes": "monday"}, expected: false},
		{item: &journeyWithEmptyDtArr, conditions: map[string]string{"dayTypes": "monday"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"departureTime": "00:00"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"arrivalTime": "00:00"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"firstStopPointId": "11111"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"lastStopPointId": "11111"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"stopPointId": "11111"}, expected: false},
		{item: &journeyWithEmptyCallArr, conditions: map[string]string{"firstStopPointId": "11111"}, expected: false},
		{item: &journeyWithEmptyCallArr, conditions: map[string]string{"lastStopPointId": "11111"}, expected: false},
		{item: &journeyWithEmptyCallArr, conditions: map[string]string{"stopPointId": "11111"}, expected: false},
		{item: &emptyJourney, conditions: map[string]string{"gtfsTripId": "11111"}, expected: false},
		{item: &invalidJourneyLowerSide, conditions: nil, expected: false},
		{item: &invalidJourneyUpperSide, conditions: nil, expected: false},
		{item: &validJourney, conditions: nil, expected: true},
		{item: &validJourney, conditions: map[string]string{"lineId": "1A"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"routeId": "123"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"journeyPatternId": "123"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"dayTypes": "monday,tuesday"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"departureTime": "01:00"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"arrivalTime": "02:00"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"firstStopPointId": "1"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"stopPointId": "2"}, expected: true},
		{item: &validJourney, conditions: map[string]string{"lastStopPointId": "3"}, expected: true},
	}
	for _, tc := range testCases {
		matches := journeyMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestJourneysRoutes(t *testing.T) {
	testCases := []struct {
		target string
		items  []Journey
	}{
		{"/journeys",
			[]Journey{getJourneyMap()["7020205685"], getJourneyMap()["7020295685"], getJourneyMap()["7024545685"]},
		},
		{"/journeys?firstStopPointId=4600&lastStopPointId=8149&stopPointId=8171",
			[]Journey{getJourneyMap()["7020295685"]},
		},
		{"/journeys?journeyPatternId=047b0afc973ee2fd4fe92b128c3a932a",
			[]Journey{getJourneyMap()["7020205685"]},
		},
		{"/journeys?lineId=1",
			[]Journey{getJourneyMap()["7020205685"]},
		},
		{"/journeys?routeId=1501146007035",
			[]Journey{getJourneyMap()["7020295685"]},
		},
		{"/journeys?dayTypes=monday,tuesday",
			[]Journey{getJourneyMap()["7020205685"], getJourneyMap()["7020295685"], getJourneyMap()["7024545685"]},
		},
		{"/journeys?departureTime=14:43:00&arrivalTime=14:44:45",
			[]Journey{getJourneyMap()["7020205685"]},
		},
		{"/journeys?departureTime=14:43:00",
			[]Journey{getJourneyMap()["7020205685"]},
		},
		{"/journeys?arrivalTime=14:44:45",
			[]Journey{getJourneyMap()["7020205685"]},
		},
		{"/journeys?gtfsTripId=7020295685",
			[]Journey{getJourneyMap()["7020295685"]},
		},
		{"/journeys/7020295685",
			[]Journey{getJourneyMap()["7020295685"]},
		},
		{"/journeys/1A_0630_8149_4600",
			[]Journey{getJourneyMap()["7020295685"]},
		},
		{"/journeys/foobar",
			[]Journey{},
		},
		{"/journeys?dayTypes=saturday",
			[]Journey{},
		},
		{"/journeys?firstStopPointId=foo&lastStopPointId=bar&stopPointId=baz",
			[]Journey{},
		},
		{"/journeys?firstStopPointId=baz",
			[]Journey{},
		},
		{"/journeys?lastStopPointId=baz",
			[]Journey{},
		},
		{"/journeys?stopPointId=baz",
			[]Journey{},
		},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectJourneyRoutes(r, ctx)

		gotResponse := getJourneySuccessResponse(t, r, w, tc.target)
		expectedResponse := journeySuccessResponse{
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

func getJourneyMap() map[string]Journey {
	result := make(map[string]Journey)

	journeys := []struct {
		id                   string
		line                 string
		activityId           string
		route                string
		journeyPattern       string
		departureTime        string
		arrivalTime          string
		headSign             string
		directionId          string
		wheelchairAccessible bool
		gtfs                 JourneyGtfsInfo
		dayTypes             []string
		dayTypeExceptions    []DayTypeException
		calls                []JourneyCall
	}{
		{
			"111111111",
			"-1",
			"-1_0720_7017_3615",
			"111111111",
			"c01c71b0c9f456ba21f498a1dca54b3b",
			"07:20:00",
			"07:21:00",
			"Foobar",
			"0",
			false,
			JourneyGtfsInfo{TripId: "111111111"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{},
			[]JourneyCall{
				{"07:20:00", "07:20:00", getStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getStopPointMap()["7017"]},
			},
		},
		{
			"7020205685",
			"1",
			"1_1443_7015_7017",
			"1504270174600",
			"047b0afc973ee2fd4fe92b128c3a932a",
			"14:43:00",
			"14:44:45",
			"Vatiala",
			"1",
			false,
			JourneyGtfsInfo{TripId: "7020205685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"14:43:00", "14:43:00", getStopPointMap()["7017"]},
				{"14:44:45", "14:44:45", getStopPointMap()["7015"]},
			},
		},
		{
			"7020295685",
			"1A",
			"1A_0630_8149_4600",
			"1501146007035",
			"9bc7403ad27267edbfbd63c3e92e5afa",
			"06:30:00",
			"06:32:30",
			"Lentoasema",
			"0",
			false,
			JourneyGtfsInfo{TripId: "7020295685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"06:30:00", "06:30:00", getStopPointMap()["4600"]},
				{"06:31:30", "06:31:30", getStopPointMap()["8171"]},
				{"06:32:30", "06:32:30", getStopPointMap()["8149"]},
			},
		},
		{
			"7024545685",
			"3A",
			"3A_0720_3607_3615",
			"1517136151028",
			"65f51d2f85284af2fad1305c0ce71033",
			"07:20:00",
			"07:21:00",
			"Lentävänniemi",
			"0",
			false,
			JourneyGtfsInfo{TripId: "7024545685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"07:20:00", "07:20:00", getStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getStopPointMap()["3607"]},
			},
		},
	}

	for _, tc := range journeys {
		result[tc.id] = Journey{
			Url:                  journeyUrl(tc.id),
			ActivityUrl:          journeyActivityUrl(tc.activityId),
			LineUrl:              lineUrl(tc.line),
			RouteUrl:             routeUrl(tc.route),
			JourneyPatternUrl:    journeyPatternUrl(tc.journeyPattern),
			DepartureTime:        tc.departureTime,
			ArrivalTime:          tc.arrivalTime,
			HeadSign:             tc.headSign,
			Direction:            tc.directionId,
			WheelchairAccessible: tc.wheelchairAccessible,
			GtfsInfo:             tc.gtfs,
			DayTypes:             tc.dayTypes,
			DayTypeExceptions:    tc.dayTypeExceptions,
			Calls:                tc.calls,
		}
	}

	return result
}

type journeySuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Journey      `json:"body"`
}

func getJourneySuccessResponse(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) journeySuccessResponse {
	serveHttp(t, r, w, target)

	var response journeySuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}
