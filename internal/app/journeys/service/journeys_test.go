package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/testutil"
	"testing"
	"time"
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
		id         string
		item       *model.Journey
		conditions map[string]string
		expected   bool
	}{
		{id: "1", item: &emptyJourney, conditions: map[string]string{"lineId": "1"}, expected: false},
		{id: "2", item: &emptyJourney, conditions: map[string]string{"routeId": "1"}, expected: false},
		{id: "3", item: &emptyJourney, conditions: map[string]string{"journeyPatternId": "1"}, expected: false},
		{id: "4", item: &emptyJourney, conditions: map[string]string{"dayTypes": "monday"}, expected: false},
		{id: "5", item: &journeyWithEmptyDtArr, conditions: map[string]string{"dayTypes": "monday"}, expected: false},
		{id: "6", item: &emptyJourney, conditions: map[string]string{"departureTime": "00:00"}, expected: false},
		{id: "7", item: &emptyJourney, conditions: map[string]string{"arrivalTime": "00:00"}, expected: false},
		{id: "8", item: &emptyJourney, conditions: map[string]string{"firstStopPointId": "11111"}, expected: false},
		{id: "9", item: &emptyJourney, conditions: map[string]string{"lastStopPointId": "11111"}, expected: false},
		{id: "10", item: &emptyJourney, conditions: map[string]string{"stopPointId": "11111"}, expected: false},
		{id: "11", item: &journeyWithEmptyCallArr, conditions: map[string]string{"firstStopPointId": "11111"}, expected: false},
		{id: "12", item: &journeyWithEmptyCallArr, conditions: map[string]string{"lastStopPointId": "11111"}, expected: false},
		{id: "13", item: &journeyWithEmptyCallArr, conditions: map[string]string{"stopPointId": "11111"}, expected: false},
		{id: "14", item: &emptyJourney, conditions: map[string]string{"gtfsTripId": "11111"}, expected: false},
		{id: "15", item: &invalidJourneyLowerSide, conditions: nil, expected: false},
		{id: "16", item: &invalidJourneyUpperSide, conditions: nil, expected: false},
		{id: "17", item: &validJourney, conditions: nil, expected: true},
		{id: "18", item: &validJourney, conditions: map[string]string{"lineId": "1A"}, expected: true},
		{id: "19", item: &validJourney, conditions: map[string]string{"routeId": "123"}, expected: true},
		{id: "20", item: &validJourney, conditions: map[string]string{"journeyPatternId": "123"}, expected: true},
		{id: "21", item: &validJourney, conditions: map[string]string{"dayTypes": "monday,tuesday"}, expected: true},
		{id: "22", item: &validJourney, conditions: map[string]string{"departureTime": "01:00"}, expected: true},
		{id: "23", item: &validJourney, conditions: map[string]string{"arrivalTime": "02:00"}, expected: true},
		{id: "24", item: &validJourney, conditions: map[string]string{"firstStopPointId": "1"}, expected: true},
		{id: "25", item: &validJourney, conditions: map[string]string{"stopPointId": "2"}, expected: true},
		{id: "26", item: &validJourney, conditions: map[string]string{"lastStopPointId": "3"}, expected: true},
	}
	for _, tc := range testCases {
		matches := journeyMatchesConditions(tc.item, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}

func TestJourneysService_Search(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Journeys: &repository.JourneysJourneyRepository{
			All: []*model.Journey{
				{Id: "1", ValidFrom: "2023-01-01", ValidTo: "2100-12-31", Line: &model.Line{Name: "1"}, Route: &model.Route{Id: "route1"}, JourneyPattern: &model.JourneyPattern{Id: "pattern1"}, DayTypes: []string{"monday"}, DepartureTime: "01:00", ArrivalTime: "02:00", Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}, {StopPoint: &model.StopPoint{ShortName: "SP2"}}}},
				{Id: "2", ValidFrom: "2023-01-01", ValidTo: "2100-12-31", Line: &model.Line{Name: "2"}, Route: &model.Route{Id: "route2"}, JourneyPattern: &model.JourneyPattern{Id: "pattern2"}, DayTypes: []string{"tuesday"}, DepartureTime: "03:00", ArrivalTime: "04:00", Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP3"}}, {StopPoint: &model.StopPoint{ShortName: "SP4"}}}},
			},
		},
	}
	service := JourneysService{Repository: dataStore}

	testCases := []struct {
		id       string
		params   map[string]string
		expected int
	}{
		{"1", map[string]string{"lineId": "1"}, 1},
		{"2", map[string]string{"routeId": "route1"}, 1},
		{"3", map[string]string{"journeyPatternId": "pattern1"}, 1},
		{"4", map[string]string{"dayTypes": "monday"}, 1},
		{"5", map[string]string{"departureTime": "01:00"}, 1},
		{"6", map[string]string{"arrivalTime": "02:00"}, 1},
		{"7", map[string]string{"firstStopPointId": "SP1"}, 1},
		{"8", map[string]string{"lastStopPointId": "SP2"}, 1},
		{"9", map[string]string{"stopPointId": "SP3"}, 1},
		{"10", map[string]string{"lineId": "NonExistent"}, 0},
	}

	for _, tc := range testCases {
		result := service.Search(tc.params)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, len(result), tc.id)
	}
}

func TestJourneysService_GetOneById(t *testing.T) {
	dataStore := &repository.JourneysRepository{
		Journeys: &repository.JourneysJourneyRepository{
			ById: map[string]*model.Journey{
				"1": {Id: "1"},
				"2": {Id: "2"},
			},
			ByActivityId: map[string]*model.Journey{
				"activity1": {Id: "1"},
				"activity2": {Id: "2"},
			},
		},
	}
	service := JourneysService{Repository: dataStore}

	testCases := []struct {
		id       string
		expected *model.Journey
		err      error
	}{
		{"1", &model.Journey{Id: "1"}, nil},
		{"2", &model.Journey{Id: "2"}, nil},
		{"activity1", &model.Journey{Id: "1"}, nil},
		{"activity2", &model.Journey{Id: "2"}, nil},
		{"3", nil, model.ErrNoSuchElement},
	}

	for _, tc := range testCases {
		result, err := service.GetOneById(tc.id)
		if err != nil && tc.err == nil {
			t.Error(err)
		} else if err != nil && tc.err == nil {
			testutil.CompareVariablesAndPrintResults(t, tc.err, err, tc.id)
		} else {
			testutil.CompareVariablesAndPrintResults(t, tc.expected, result, tc.id)
		}
	}
}

func TestJourneyMatchesConditions2(t *testing.T) {
	now := time.Now()
	curDay := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

	testCases := []struct {
		id         string
		journey    *model.Journey
		conditions map[string]string
		expected   bool
	}{
		{"1", nil, nil, false},
		{"2", &model.Journey{ValidFrom: curDay, ValidTo: curDay}, nil, true},
		{"3", &model.Journey{ValidFrom: "2022-01-01", ValidTo: "2022-12-31"}, nil, false},
		{"4", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Line: &model.Line{Name: "1"}}, map[string]string{"lineId": "1"}, true},
		{"5", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Line: &model.Line{Name: "2"}}, map[string]string{"lineId": "1"}, false},
		{"6", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Route: &model.Route{Id: "route1"}}, map[string]string{"routeId": "route1"}, true},
		{"7", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Route: &model.Route{Id: "route2"}}, map[string]string{"routeId": "route1"}, false},
		{"8", &model.Journey{ValidFrom: curDay, ValidTo: curDay, JourneyPattern: &model.JourneyPattern{Id: "pattern1"}}, map[string]string{"journeyPatternId": "pattern1"}, true},
		{"9", &model.Journey{ValidFrom: curDay, ValidTo: curDay, JourneyPattern: &model.JourneyPattern{Id: "pattern2"}}, map[string]string{"journeyPatternId": "pattern1"}, false},
		{"10", &model.Journey{ValidFrom: curDay, ValidTo: curDay, DayTypes: []string{"monday"}}, map[string]string{"dayTypes": "monday"}, true},
		{"11", &model.Journey{ValidFrom: curDay, ValidTo: curDay, DayTypes: []string{"tuesday"}}, map[string]string{"dayTypes": "monday"}, false},
		{"12", &model.Journey{ValidFrom: curDay, ValidTo: curDay, DepartureTime: "01:00"}, map[string]string{"departureTime": "01:00"}, true},
		{"13", &model.Journey{ValidFrom: curDay, ValidTo: curDay, DepartureTime: "02:00"}, map[string]string{"departureTime": "01:00"}, false},
		{"14", &model.Journey{ValidFrom: curDay, ValidTo: curDay, ArrivalTime: "02:00"}, map[string]string{"arrivalTime": "02:00"}, true},
		{"15", &model.Journey{ValidFrom: curDay, ValidTo: curDay, ArrivalTime: "03:00"}, map[string]string{"arrivalTime": "02:00"}, false},
		{"16", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}}}, map[string]string{"firstStopPointId": "SP1"}, true},
		{"17", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP2"}}}}, map[string]string{"firstStopPointId": "SP1"}, false},
		{"18", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}, {StopPoint: &model.StopPoint{ShortName: "SP2"}}}}, map[string]string{"lastStopPointId": "SP2"}, true},
		{"19", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}, {StopPoint: &model.StopPoint{ShortName: "SP3"}}}}, map[string]string{"lastStopPointId": "SP2"}, false},
		{"20", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}, {StopPoint: &model.StopPoint{ShortName: "SP2"}}}}, map[string]string{"stopPointId": "SP2"}, true},
		{"21", &model.Journey{ValidFrom: curDay, ValidTo: curDay, Calls: []*model.JourneyCall{{StopPoint: &model.StopPoint{ShortName: "SP1"}}, {StopPoint: &model.StopPoint{ShortName: "SP3"}}}}, map[string]string{"stopPointId": "SP2"}, false},
	}

	for _, tc := range testCases {
		matches := journeyMatchesConditions(tc.journey, tc.conditions)
		testutil.CompareVariablesAndPrintResults(t, tc.expected, matches, tc.id)
	}
}
