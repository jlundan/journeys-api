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

		response := getJourneySuccessResponse(t, r, w, tc.target)

		dataSize := len(tc.items)
		if success := validateCommonResponseFields(t, response.Status, response.Data, uint16(dataSize)); !success {
			break
		}
		if len(response.Body) != dataSize {
			t.Errorf("expected %v, got %v", dataSize, len(response.Body))
			break
		}
		for i, l := range response.Body {
			var diffs []FieldDiff
			var err = compareVariables(tc.items[i], l, tc.target, &diffs)
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
