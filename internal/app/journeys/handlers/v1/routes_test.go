//go:build journeys_routes_tests || journeys_tests || all_tests

package v1

import (
	"fmt"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	dataService := newJourneysTestDataService(t)

	one := handlerConfig{handler: HandleGetOneRoute(dataService, ""), url: "/v1/routes/{name}"}
	all := handlerConfig{handler: HandleGetAllRoutes(dataService, ""), url: "/v1/routes"}

	rm := getRouteMap()

	testCases := []routerTestCase[Route]{
		{"/v1/routes",
			[]Route{rm[routeUrl("111111111")], rm[routeUrl("1501146007035")], rm[routeUrl("1504270174600")], rm[routeUrl("1517136151028")]},
			false, all,
		},
		{"/v1/routes/1501146007035",
			[]Route{rm[routeUrl("1501146007035")]},
			false, one,
		},
		{"/v1/routes?lineId=1A",
			[]Route{rm[routeUrl("1501146007035")]},
			false, all,
		},
		{"/v1/routes/foobar",
			[]Route{},
			false, one,
		},
	}

	runRouterTestCases(t, testCases)

}

func getRouteMap() map[string]Route {
	result := make(map[string]Route)

	routes := []Route{
		{
			Url:        routeUrl("111111111"),
			LineUrl:    routeLineUrl("-1"),
			Name:       "Näyttelijänkatu - Suupantori",
			Projection: "6144444,2387262:17,-14",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             routeJourneyPatternUrl("c01c71b0c9f456ba21f498a1dca54b3b"),
					OriginStop:      stopPointUrl("3615"),
					DestinationStop: stopPointUrl("7017"),
					Name:            "Näyttelijänkatu - Suupantori",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               routeJourneyUrl("111111111"),
					JourneyPatternUrl: routeJourneyPatternUrl("c01c71b0c9f456ba21f498a1dca54b3b"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					DayTypes:          []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
					DayTypeExceptions: []DayTypeException{},
				},
			},
		},
		{
			Url:        routeUrl("1501146007035"),
			LineUrl:    routeLineUrl("1A"),
			Name:       "Vatiala - Sudenkorennontie",
			Projection: "6147590,2397764:-104,32",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             routeJourneyPatternUrl("9bc7403ad27267edbfbd63c3e92e5afa"),
					OriginStop:      stopPointUrl("4600"),
					DestinationStop: stopPointUrl("8149"),
					Name:            "Vatiala - Sudenkorennontie",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               routeJourneyUrl("7020295685"),
					JourneyPatternUrl: routeJourneyPatternUrl("9bc7403ad27267edbfbd63c3e92e5afa"),
					DepartureTime:     "06:30:00",
					ArrivalTime:       "06:32:30",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
		{
			Url:        routeUrl("1504270174600"),
			LineUrl:    routeLineUrl("1"),
			Name:       "Suupantori - Pirkkala",
			Projection: "6146557,2364305:-19,-148:-17,-191",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             routeJourneyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					OriginStop:      stopPointUrl("7017"),
					DestinationStop: stopPointUrl("7015"),
					Name:            "Suupantori - Pirkkala",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               routeJourneyUrl("7020205685"),
					JourneyPatternUrl: routeJourneyPatternUrl("047b0afc973ee2fd4fe92b128c3a932a"),
					DepartureTime:     "14:43:00",
					ArrivalTime:       "14:44:45",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
		{
			Url:        routeUrl("1517136151028"),
			LineUrl:    routeLineUrl("3A"),
			Name:       "Näyttelijänkatu - Lavastajanpolku",
			Projection: "6144444,2387262:17,-14",
			JourneyPatterns: []RouteJourneyPattern{
				{
					Url:             routeJourneyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					OriginStop:      stopPointUrl("3615"),
					DestinationStop: stopPointUrl("3607"),
					Name:            "Näyttelijänkatu - Lavastajanpolku",
				},
			},
			Journeys: []RouteJourney{
				{
					Url:               routeJourneyUrl("7024545685"),
					JourneyPatternUrl: routeJourneyPatternUrl("65f51d2f85284af2fad1305c0ce71033"),
					DepartureTime:     "07:20:00",
					ArrivalTime:       "07:21:00",
					DayTypes: []string{
						"monday",
						"tuesday",
						"wednesday",
						"thursday",
						"friday",
					},
					DayTypeExceptions: []DayTypeException{
						{
							From: "2021-04-05",
							To:   "2021-04-05",
							Runs: "yes",
						},
						{
							From: "2021-05-13",
							To:   "2021-05-13",
							Runs: "no",
						},
					},
				},
			},
		},
	}

	for _, route := range routes {
		result[route.Url] = route
	}

	return result
}

func routeUrl(name string) string {
	return fmt.Sprintf("%v/routes/%v", "", name)
}

func routeLineUrl(name string) string {
	return fmt.Sprintf("%v/lines/%v", "", name)
}

func routeJourneyUrl(name string) string {
	return fmt.Sprintf("%v/journeys/%v", "", name)
}

func routeJourneyPatternUrl(name string) string {
	return fmt.Sprintf("%v/journey-patterns/%v", "", name)
}
