//go:build journeys_stops_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestStopPointRoutes(t *testing.T) {
	dataService := newJourneysTestDataService()

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
