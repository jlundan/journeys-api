package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
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
