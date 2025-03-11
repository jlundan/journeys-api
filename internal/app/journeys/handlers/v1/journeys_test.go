//go:build journeys_journeys_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestJourneysRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneJourney(dataService, "", ""), url: "/v1/journeys/{name}"}
	all := handlerConfig{handler: HandleGetAllJourneys(dataService, "", ""), url: "/v1/journeys"}

	jm := getJourneyMap()
	testCases := []routerTestCase[Journey]{
		{"/v1/journeys",
			[]Journey{jm["7020205685"], jm["7020295685"], jm["7024545685"]}, false, all,
		},
		{"/v1/journeys?firstStopPointId=4600&lastStopPointId=8149&stopPointId=8171",
			[]Journey{jm["7020295685"]}, false, all,
		},
		{"/v1/journeys?journeyPatternId=047b0afc973ee2fd4fe92b128c3a932a",
			[]Journey{jm["7020205685"]}, false, all,
		},
		{"/v1/journeys?lineId=1",
			[]Journey{jm["7020205685"]}, false, all,
		},
		{"/v1/journeys?routeId=1501146007035",
			[]Journey{jm["7020295685"]}, false, all,
		},
		{"/v1/journeys?dayTypes=monday,tuesday",
			[]Journey{jm["7020205685"], jm["7020295685"], jm["7024545685"]}, false, all,
		},
		{"/v1/journeys?departureTime=14:43:00&arrivalTime=14:44:45",
			[]Journey{jm["7020205685"]}, false, all,
		},
		{"/v1/journeys?departureTime=14:43:00",
			[]Journey{jm["7020205685"]}, false, all,
		},
		{"/v1/journeys?arrivalTime=14:44:45",
			[]Journey{jm["7020205685"]}, false, all,
		},
		{"/v1/journeys?gtfsTripId=7020295685",
			[]Journey{jm["7020295685"]}, false, all,
		},
		{"/v1/journeys/7020295685",
			[]Journey{jm["7020295685"]}, false, one,
		},
		{"/v1/journeys/1A_0630_8149_4600",
			[]Journey{jm["7020295685"]}, false, one,
		},
		{"/v1/journeys/foobar",
			[]Journey{}, false, one,
		},
		{"/v1/journeys?dayTypes=saturday",
			[]Journey{}, false, all,
		},
		{"/v1/journeys?firstStopPointId=foo&lastStopPointId=bar&stopPointId=baz",
			[]Journey{}, false, all,
		},
		{"/v1/journeys?firstStopPointId=baz",
			[]Journey{}, false, all,
		},
		{"/v1/journeys?lastStopPointId=baz",
			[]Journey{}, false, all,
		},
		{"/v1/journeys?stopPointId=baz",
			[]Journey{}, false, all,
		},
	}

	runRouterTestCases(t, testCases)
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
				{"07:20:00", "07:20:00", getJourneyStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getJourneyStopPointMap()["7017"]},
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
				{"14:43:00", "14:43:00", getJourneyStopPointMap()["7017"]},
				{"14:44:45", "14:44:45", getJourneyStopPointMap()["7015"]},
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
				{"06:30:00", "06:30:00", getJourneyStopPointMap()["4600"]},
				{"06:31:30", "06:31:30", getJourneyStopPointMap()["8171"]},
				{"06:32:30", "06:32:30", getJourneyStopPointMap()["8149"]},
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
				{"07:20:00", "07:20:00", getJourneyStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getJourneyStopPointMap()["3607"]},
			},
		},
	}

	for _, tc := range journeys {
		result[tc.id] = Journey{
			Url:                  journeyUrl(tc.id),
			ActivityUrl:          journeyActivityUrl(tc.activityId),
			LineUrl:              journeyLineUrl(tc.line),
			RouteUrl:             journeyRouteUrl(tc.route),
			JourneyPatternUrl:    journeyJourneyPatternUrl(tc.journeyPattern),
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

func getJourneyStopPointMap() map[string]JourneyStopPoint {
	result := make(map[string]JourneyStopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality JourneyMunicipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getJourneyMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getJourneyMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getJourneyMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getJourneyMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getJourneyMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getJourneyMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getJourneyMunicipalityMap()["837"]},
	}

	for _, tc := range stopPoints {
		result[tc.id] = JourneyStopPoint{
			Url:          journeyStopPointUrl(tc.id),
			ShortName:    tc.id,
			Name:         tc.name,
			Location:     tc.location,
			TariffZone:   tc.tariffZone,
			Municipality: tc.municipality,
		}
	}

	return result
}

func getJourneyMunicipalityMap() map[string]JourneyMunicipality {
	municipalities := make(map[string]JourneyMunicipality)
	municipalities["211"] = JourneyMunicipality{
		Url:       journeysMunicipalityUrl("211"),
		ShortName: "211",
		Name:      "Kangasala",
	}

	municipalities["604"] = JourneyMunicipality{
		Url:       journeysMunicipalityUrl("604"),
		ShortName: "604",
		Name:      "Pirkkala",
	}

	municipalities["837"] = JourneyMunicipality{
		Url:       journeysMunicipalityUrl("837"),
		ShortName: "837",
		Name:      "Tampere",
	}
	return municipalities
}

func journeyUrl(name string) string {
	return fmt.Sprintf("%v/journeys/%v", "", name)
}

func journeyStopPointUrl(name string) string {
	return fmt.Sprintf("%v/stop-points/%v", "", name)
}

func journeyLineUrl(name string) string {
	return fmt.Sprintf("%v/lines/%v", "", name)
}

func journeyRouteUrl(name string) string {
	return fmt.Sprintf("%v/routes/%v", "", name)
}

func journeyJourneyPatternUrl(name string) string {
	return fmt.Sprintf("%v/journey-patterns/%v", "", name)
}

func journeyActivityUrl(name string) string {
	return fmt.Sprintf("%v/vehicle-activity?journeyRef=%v", "", name)
}

func journeysMunicipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", "", name)
}
