//go:build journeys_journeypattern_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestJourneyPatternsRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneJourneyPattern(dataService, ""), url: "/v1/journey-patterns/{name}"}
	all := handlerConfig{handler: HandleGetAllJourneyPatterns(dataService, ""), url: "/v1/journey-patterns"}

	jpm := getJourneyPatternMap()

	testCases := []routerTestCase[JourneyPattern]{
		{"/v1/journey-patterns",
			[]JourneyPattern{
				jpm["047b0afc973ee2fd4fe92b128c3a932a"],
				jpm["65f51d2f85284af2fad1305c0ce71033"],
				jpm["9bc7403ad27267edbfbd63c3e92e5afa"],
				jpm["c01c71b0c9f456ba21f498a1dca54b3b"],
			}, false, all,
		},
		{"/v1/journey-patterns?lineId=1",
			[]JourneyPattern{
				jpm["047b0afc973ee2fd4fe92b128c3a932a"],
			}, false, all,
		},
		{"/v1/journey-patterns?firstStopPointId=4600&lastStopPointId=8149&stopPointId=8171",
			[]JourneyPattern{
				jpm["9bc7403ad27267edbfbd63c3e92e5afa"],
			}, false, all,
		},
		{"/v1/journey-patterns/047b0afc973ee2fd4fe92b128c3a932a?exclude-fields=name,stopPoints.municipality.url",
			[]JourneyPattern{
				{
					Url:             journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					LineUrl:         journeyPatternLineUrl("1"),
					RouteUrl:        journeyPatternRouteUrl("1504270174600"),
					OriginStop:      journeyPatternStopPointUrl("7017"),
					DestinationStop: journeyPatternStopPointUrl("7015"),
					Direction:       "1",
					StopPoints: []JourneyPatternStopPoint{
						{
							Url:        journeyPatternStopPointUrl("7017"),
							ShortName:  "7017",
							Name:       "Suupantori",
							Location:   "61.46546,23.64219",
							TariffZone: "B",
							Municipality: JourneyPatternMunicipality{
								ShortName: "604",
								Name:      "Pirkkala",
							},
						},
						{
							Url:        journeyPatternStopPointUrl("7015"),
							ShortName:  "7015",
							Name:       "Pirkkala",
							Location:   "61.4659,23.64734",
							TariffZone: "B",
							Municipality: JourneyPatternMunicipality{
								ShortName: "604",
								Name:      "Pirkkala",
							},
						},
					},
					Journeys: []JourneyPatternJourney{
						{
							Url:               journeyPatternJourneyUrl("7020205685"),
							JourneyPatternUrl: journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
							DepartureTime:     "14:43:00",
							ArrivalTime:       "14:44:45",
							HeadSign:          "Vatiala",
							DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
							DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
						},
					},
				},
			}, false, one,
		},
		{"/v1/journey-patterns/foobar",
			[]JourneyPattern{}, false, one,
		},
		{"/v1/journey-patterns?lastStopPointId=foobar",
			[]JourneyPattern{}, false, all,
		},
		{"/v1/journey-patterns?firstStopPointId=foobar",
			[]JourneyPattern{}, false, all,
		},
		{"/v1/journey-patterns?stopPointId=foobar",
			[]JourneyPattern{}, false, all,
		},
	}

	runRouterTestCases(t, testCases)
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
		stopPoints      []JourneyPatternStopPoint
		journeys        []JourneyPatternJourney
	}{
		{"047b0afc973ee2fd4fe92b128c3a932a", "1", "1504270174600", "7017",
			"7015", "Suupantori - Pirkkala", "1",
			[]JourneyPatternStopPoint{getJourneyPatternStopPointMap()["7017"], getJourneyPatternStopPointMap()["7015"]},
			[]JourneyPatternJourney{
				{
					Url:               journeyPatternJourneyUrl("7020205685"),
					JourneyPatternUrl: journeyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					DepartureTime:     "14:43:00",
					ArrivalTime:       "14:44:45",
					HeadSign:          "Vatiala",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
			}},

		{"65f51d2f85284af2fad1305c0ce71033", "3A", "1517136151028", "3615",
			"3607", "Näyttelijänkatu - Lavastajanpolku", "0",
			[]JourneyPatternStopPoint{getJourneyPatternStopPointMap()["3615"], getJourneyPatternStopPointMap()["3607"]},
			[]JourneyPatternJourney{
				{
					Url:               journeyPatternJourneyUrl("7024545685"),
					JourneyPatternUrl: journeyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					HeadSign:          "Lentävänniemi",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
				{
					Url:               journeyPatternJourneyUrl("123456789"),
					JourneyPatternUrl: journeyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					HeadSign:          "Lentävänniemi",
					DayTypes:          []string{"saturday", "sunday"},
					DayTypeExceptions: []DayTypeException{},
				},
			}},

		{"9bc7403ad27267edbfbd63c3e92e5afa", "1A", "1501146007035", "4600",
			"8149", "Vatiala - Sudenkorennontie", "0",
			[]JourneyPatternStopPoint{getJourneyPatternStopPointMap()["4600"], getJourneyPatternStopPointMap()["8171"], getJourneyPatternStopPointMap()["8149"]},
			[]JourneyPatternJourney{
				{
					Url:               journeyPatternJourneyUrl("7020295685"),
					JourneyPatternUrl: journeyPatternUrl("9bc7403ad27267edbfbd63c3e92e5afa"),
					DepartureTime:     "06:30:00",
					ArrivalTime:       "06:32:30",
					HeadSign:          "Lentoasema",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
				},
			}},

		{"c01c71b0c9f456ba21f498a1dca54b3b", "-1", "111111111", "3615",
			"7017", "Näyttelijänkatu - Suupantori", "0",
			[]JourneyPatternStopPoint{getJourneyPatternStopPointMap()["3615"], getJourneyPatternStopPointMap()["7017"]},
			[]JourneyPatternJourney{
				{
					Url:               journeyPatternJourneyUrl("111111111"),
					JourneyPatternUrl: journeyPatternUrl("c01c71b0c9f456ba21f498a1dca54b3b"),
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
			LineUrl:         journeyPatternLineUrl(tc.line),
			RouteUrl:        journeyPatternRouteUrl(tc.route),
			OriginStop:      journeyPatternStopPointUrl(tc.originStop),
			DestinationStop: journeyPatternStopPointUrl(tc.destinationStop),
			Direction:       tc.direction,
			Name:            tc.name,
			StopPoints:      tc.stopPoints,
			Journeys:        tc.journeys,
		}
	}

	return result
}

func getJourneyPatternStopPointMap() map[string]JourneyPatternStopPoint {
	result := make(map[string]JourneyPatternStopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality JourneyPatternMunicipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getJourneyPatternStopPointMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getJourneyPatternStopPointMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getJourneyPatternStopPointMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getJourneyPatternStopPointMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getJourneyPatternStopPointMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getJourneyPatternStopPointMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getJourneyPatternStopPointMunicipalityMap()["837"]},
	}

	for _, tc := range stopPoints {
		result[tc.id] = JourneyPatternStopPoint{
			Url:          journeyPatternStopPointUrl(tc.id),
			ShortName:    tc.id,
			Name:         tc.name,
			Location:     tc.location,
			TariffZone:   tc.tariffZone,
			Municipality: tc.municipality,
		}
	}

	return result
}

func getJourneyPatternStopPointMunicipalityMap() map[string]JourneyPatternMunicipality {
	municipalities := make(map[string]JourneyPatternMunicipality)
	municipalities["211"] = JourneyPatternMunicipality{
		Url:       journeyPatternMunicipalityUrl("211"),
		ShortName: "211",
		Name:      "Kangasala",
	}

	municipalities["604"] = JourneyPatternMunicipality{
		Url:       journeyPatternMunicipalityUrl("604"),
		ShortName: "604",
		Name:      "Pirkkala",
	}

	municipalities["837"] = JourneyPatternMunicipality{
		Url:       journeyPatternMunicipalityUrl("837"),
		ShortName: "837",
		Name:      "Tampere",
	}
	return municipalities
}

func journeyPatternUrl(name string) string {
	return fmt.Sprintf("%v/journey-patterns/%v", "", name)
}

func journeyPatternStopPointUrl(name string) string {
	return fmt.Sprintf("%v/stop-points/%v", "", name)
}

func journeyPatternLineUrl(name string) string {
	return fmt.Sprintf("%v/lines/%v", "", name)
}

func journeyPatternRouteUrl(name string) string {
	return fmt.Sprintf("%v/routes/%v", "", name)
}

func journeyPatternMunicipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", "", name)
}

func journeyPatternJourneyUrl(name string) string {
	return fmt.Sprintf("%v/journeys/%v", "", name)
}
