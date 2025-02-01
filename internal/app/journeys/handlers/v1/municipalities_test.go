//go:build journeys_municipalities_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestMunicipalitiesRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneMunicipality(dataService, ""), url: "/v1/municipalities/{name}"}
	all := handlerConfig{handler: HandleGetAllMunicipalities(dataService, ""), url: "/v1/municipalities"}

	testCases := []routerTestCase[Municipality]{
		{"/v1/municipalities", []Municipality{
			{municipalityUrl("211"), "211", "Kangasala"},
			{municipalityUrl("604"), "604", "Pirkkala"},
			{municipalityUrl("837"), "837", "Tampere"},
			{municipalityUrl("980"), "980", "Ylöjärvi"}},
			false, all,
		},
		{"/v1/municipalities?name=Pirkkala", []Municipality{
			{municipalityUrl("604"), "604", "Pirkkala"},
		}, false, all},
		{"/v1/municipalities?name=Tampere&shortName=837", []Municipality{
			{municipalityUrl("837"), "837", "Tampere"},
		}, false, all},
		{"/v1/municipalities?shortName=837&exclude-fields=name,shortName", []Municipality{
			{municipalityUrl("837"), "", ""},
		}, false, all},
		{"/v1/municipalities/837", []Municipality{
			{municipalityUrl("837"), "837", "Tampere"},
		}, false, one},
		{"/v1/municipalities/foobar", []Municipality{}, false, one},
	}

	runRouterTestCases(t, testCases)
}

func municipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", "", name)
}
