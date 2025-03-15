//go:build journeys_stops_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
	"time"
)

func TestStopPointRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneStopPoint(dataService, ""), url: "/v1/stop-points/{name}"}
	all := handlerConfig{handler: HandleGetAllStopPoints(dataService, ""), url: "/v1/stop-points"}

	var sp = getStopPointMap()
	testCases := []routerTestCase[StopPoint]{
		{"/v1/stop-points",
			[]StopPoint{sp["3607"], sp["3615"], sp["4600"], sp["7015"], sp["7017"], sp["8149"], sp["8171"]},
			false, all,
		},
		{"/v1/stop-points?tariffZone=B",
			[]StopPoint{sp["3607"], sp["3615"], sp["4600"], sp["7015"], sp["7017"], sp["8171"]},
			false, all,
		},
		{"/v1/stop-points?municipalityName=Tampere",
			[]StopPoint{sp["3607"], sp["3615"]},
			false, all,
		},
		{"/v1/stop-points?municipalityShortName=837",
			[]StopPoint{sp["3607"], sp["3615"]}, false, all,
		},
		{"/v1/stop-points?location=61.4445,23.87235",
			[]StopPoint{sp["3615"]}, false, all,
		},
		{"/v1/stop-points?location=61,23:62,23.8",
			[]StopPoint{sp["7015"], sp["7017"]}, false, all,
		},
		{"/v1/stop-points?name=Pirkkala",
			[]StopPoint{sp["7015"]}, false, all,
		},
		{"/v1/stop-points?name=Sudenkorennontie&shortName=81",
			[]StopPoint{sp["8149"]}, false, all,
		},
		{"/v1/stop-points?shortName=4600&exclude-fields=name,shortName,municipality.shortName", []StopPoint{
			{
				stopPointUrl("4600"),
				"",
				"",
				"61.47561,23.97756",
				"B",
				StopPointMunicipality{
					Url:  stopPointMunicipalityUrl("211"),
					Name: "Kangasala",
				},
			},
		}, false, all},
		{"/v1/stop-points/4600", []StopPoint{sp["4600"]}, false, one},
		{"/v1/stop-points/foobar", []StopPoint{}, false, one},
		{"/v1/stop-points?location=6123:62,23.8", []StopPoint{}, false, all},
		{"/v1/stop-points?location=61,23:6223.8", []StopPoint{}, false, all},
		{"/v1/stop-points?location=A1,23:62,23.8", []StopPoint{}, false, all},
		{"/v1/stop-points?location=61,23:62,23.X", []StopPoint{}, false, all},
		{"/v1/stop-points?location=61,2X:62,23.8", []StopPoint{}, false, all},
		{"/v1/stop-points?location=61,23:6X,23.8", []StopPoint{}, false, all},
	}

	runRouterTestCases(t, testCases)
}

func getStopPointMap() map[string]StopPoint {
	result := make(map[string]StopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality StopPointMunicipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getStopPointMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getStopPointMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getStopPointMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getStopPointMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getStopPointMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getStopPointMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getStopPointMunicipalityMap()["837"]},
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

func getStopPointMunicipalityMap() map[string]StopPointMunicipality {
	municipalities := make(map[string]StopPointMunicipality)
	municipalities["211"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("211"),
		ShortName: "211",
		Name:      "Kangasala",
	}

	municipalities["604"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("604"),
		ShortName: "604",
		Name:      "Pirkkala",
	}

	municipalities["837"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("837"),
		ShortName: "837",
		Name:      "Tampere",
	}
	return municipalities
}

func stopPointUrl(name string) string {
	return fmt.Sprintf("%v/stop-points/%v", "", name)
}
func stopPointMunicipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", "", name)
}

func TestStopPointJourneyRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	journeys := handlerConfig{handler: HandleGetJourneysForStopPoint(dataService, "", "", false), url: "/v1/stop-points/{name}/journeys"}
	activeJourneys := handlerConfig{handler: HandleGetJourneysForStopPoint(dataService, "", "", true), url: "/v1/stop-points/{name}/journeys/active"}

	// Get sample journeys for specific stop points
	stopPoint3607Journeys := getJourneysForStopPoint("3607")
	stopPoint7015Journeys := getJourneysForStopPoint("7015")

	testCases := []routerTestCase[StopPointJourney]{
		{"/v1/stop-points/3607/journeys", stopPoint3607Journeys, false, journeys},
		{"/v1/stop-points/7015/journeys", stopPoint7015Journeys, false, journeys},
		{"/v1/stop-points/3607/journeys/active", filterActiveJourneys(stopPoint3607Journeys), false, activeJourneys},
		{"/v1/stop-points/7015/journeys/active", filterActiveJourneys(stopPoint7015Journeys), false, activeJourneys},
		{"/v1/stop-points/nonexistent/journeys", []StopPointJourney{}, false, journeys},
		{"/v1/stop-points/nonexistent/journeys/active", []StopPointJourney{}, false, activeJourneys},
	}

	runRouterTestCases(t, testCases)
}

// Helper function to get journeys for a specific stop point
func getJourneysForStopPoint(stopPointId string) []StopPointJourney {
	// This would normally be populated from your test data
	// For this example I'll create sample data based on the structure in convertStopPointJourney

	if stopPointId == "3607" {
		return []StopPointJourney{
			{
				JourneyUrl:           "/journeys/123456789",
				StopPointUrl:         "/stop-points/3607",
				ActivityUrl:          "/vehicle-activity?journeyRef=3A_0720_3607_3615",
				LineUrl:              "/lines/3A",
				RouteUrl:             "/routes/1517136151028",
				JourneyPatternUrl:    "/journey-patterns/65f51d2f85284af2fad1305c0ce71033",
				LineId:               "3A",
				DepartureTime:        "07:21:00",
				ArrivalTime:          "07:21:00",
				HeadSign:             "Lentävänniemi",
				Direction:            "0",
				WheelchairAccessible: false,
				GtfsInfo:             StopPointJourneyGtfsInfo{TripId: "123456789"},
				DayTypes:             []string{"saturday", "sunday"},
				DayTypeExceptions:    []StopPointDayTypeException{},
				ActiveFrom:           "2000-01-01",
				ActiveTo:             "2000-01-02",
			},
			{
				JourneyUrl:           "/journeys/7024545685",
				StopPointUrl:         "/stop-points/3607",
				ActivityUrl:          "/vehicle-activity?journeyRef=3A_0720_3607_3615",
				LineUrl:              "/lines/3A",
				RouteUrl:             "/routes/1517136151028",
				JourneyPatternUrl:    "/journey-patterns/65f51d2f85284af2fad1305c0ce71033",
				LineId:               "3A",
				DepartureTime:        "07:21:00",
				ArrivalTime:          "07:21:00",
				HeadSign:             "Lentävänniemi",
				Direction:            "0",
				WheelchairAccessible: false,
				GtfsInfo:             StopPointJourneyGtfsInfo{TripId: "7024545685"},
				DayTypes:             []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
				DayTypeExceptions: []StopPointDayTypeException{
					{From: "2021-04-05", To: "2021-04-05", Runs: "yes"},
					{From: "2021-05-13", To: "2021-05-13", Runs: "no"},
				},
				ActiveFrom: "2000-01-01",
				ActiveTo:   "2099-01-01",
			},
		}
	} else if stopPointId == "7015" {
		return []StopPointJourney{
			{
				JourneyUrl:           "/journeys/7020205685",
				StopPointUrl:         "/stop-points/7015",
				ActivityUrl:          "/vehicle-activity?journeyRef=1_1443_7015_7017",
				LineUrl:              "/lines/1",
				RouteUrl:             "/routes/1504270174600",
				JourneyPatternUrl:    "/journey-patterns/047b0afc973ee2fd4fe92b128c3a932a",
				LineId:               "1",
				DepartureTime:        "14:44:45",
				ArrivalTime:          "14:44:45",
				HeadSign:             "Vatiala",
				Direction:            "1",
				WheelchairAccessible: false,
				GtfsInfo:             StopPointJourneyGtfsInfo{TripId: "7020205685"},
				DayTypes:             []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
				DayTypeExceptions: []StopPointDayTypeException{
					{From: "2021-04-05", To: "2021-04-05", Runs: "yes"},
					{From: "2021-05-13", To: "2021-05-13", Runs: "no"},
				},
				ActiveFrom: "2000-01-01",
				ActiveTo:   "2099-01-01",
			},
		}
	}

	return []StopPointJourney{}
}

// Helper function to filter active journeys
func filterActiveJourneys(journeys []StopPointJourney) []StopPointJourney {
	now := time.Now()
	currentDate := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())

	var activeJourneys []StopPointJourney
	for _, journey := range journeys {
		if journey.ActiveFrom <= currentDate && journey.ActiveTo >= currentDate {
			activeJourneys = append(activeJourneys, journey)
		}
	}

	return activeJourneys
}
