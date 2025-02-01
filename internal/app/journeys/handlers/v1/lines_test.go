//go:build journeys_lines_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestGetLines(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneLine(dataService, ""), url: "/v1/lines/{name}"}
	all := handlerConfig{handler: HandleGetAllLines(dataService, ""), url: "/v1/lines"}

	testCases := []routerTestCase[Line]{
		{"/v1/lines", []Line{
			{lineUrl("-1"), "-1", "Foobar"},
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
			{lineUrl("3A"), "3A", "Etelä-Hervanta - Lentävänniemi"},
		}, false, all},
		{"/v1/lines?name=1", []Line{
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
		}, false, all},
		{"/v1/lines?name=1A&description=lento", []Line{
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}, false, all},
		{"/v1/lines?description=vatiala", []Line{
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}, false, all},
		{"/v1/lines?name=1&exclude-fields=name,description", []Line{
			{lineUrl("1"), "", ""},
		}, false, all},
		{"/v1/lines?name=1&exclude-fields=description", []Line{
			{lineUrl("1"), "1", ""},
		}, false, all},
		{"/v1/lines?exclude-fields=name", []Line{
			{lineUrl("-1"), "", "Foobar"},
			{lineUrl("1"), "", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "", "Vatiala - Pirkkala (lentoasema)"},
			{lineUrl("3A"), "", "Etelä-Hervanta - Lentävänniemi"},
		}, false, all},
		{"/v1/lines/1A", []Line{
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}, false, one},
		{"/v1/lines/foobar", []Line{}, false, one},
		{"/v1/lines?name=noSuchThing", []Line{}, false, all},
		{"/v1/lines?description=noSuchThing", []Line{}, false, all},
		{"/v1/lines?description=noSuchThing&name=noSuchThing", []Line{}, false, all},
		{"/v1/lines?description=vatiala&name=noSuchThing", []Line{}, false, all},
		{"/v1/lines?description=noSuchThing&name=1", []Line{}, false, all},
	}

	runRouterTestCases(t, testCases)
}

func lineUrl(name string) string {
	return fmt.Sprintf("%v/lines/%v", "", name)
}
