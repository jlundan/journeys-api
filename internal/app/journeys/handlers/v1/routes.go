package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllRoutes(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelRoutes := service.SearchRoutes(getQueryParameters(req))

		var routes []Route
		for _, ml := range modelRoutes {
			routes = append(routes, convertRoute(ml, baseUrl))
		}

		sendSuccessResponse(routes, getExcludeFieldsQueryParameter(req), rw)
	}
}

func HandleGetOneRoute(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		mr, err := service.GetOneRouteById(mux.Vars(req)["name"])
		if err != nil {
			sendSuccessResponse([]Route{}, getExcludeFieldsQueryParameter(req), rw)
			return
		}

		routes := []Route{convertRoute(mr, baseUrl)}
		sendSuccessResponse(routes, getExcludeFieldsQueryParameter(req), rw)
	}
}

func convertRoute(route *model.Route, baseUrl string) Route {
	var name string

	if len(route.JourneyPatterns) > 0 && len(route.JourneyPatterns[0].StopPoints) > 0 {
		var firstStopPoint = route.JourneyPatterns[0].StopPoints[0]
		var lastStopPoint = route.JourneyPatterns[0].StopPoints[len(route.JourneyPatterns[0].StopPoints)-1]
		name = fmt.Sprintf("%v - %v", firstStopPoint.Name, lastStopPoint.Name)
	}

	converted := Route{
		Projection: route.GeoProjection,
		Url:        fmt.Sprintf("%v%v/%v", baseUrl, routePrefix, route.Id),
		LineUrl:    fmt.Sprintf("%v%v/%v", baseUrl, linePrefix, route.Line.Name),
		Name:       name,
	}

	for _, v := range route.Journeys {
		converted.Journeys = append(converted.Journeys, RouteJourney{
			Url:               fmt.Sprintf("%v%v/%v", baseUrl, journeysPrefix, v.Id),
			JourneyPatternUrl: fmt.Sprintf("%v%v/%v", baseUrl, journeyPatternPrefix, v.JourneyPattern.Id),
			DepartureTime:     v.DepartureTime,
			ArrivalTime:       v.ArrivalTime,
			DayTypes:          v.DayTypes,
			DayTypeExceptions: makeDayTypeExceptions(v),
		})
	}

	for _, v := range route.JourneyPatterns {
		var originStop, destinationStop string

		if len(v.StopPoints) > 0 {
			originStop = fmt.Sprintf("%v%v/%v", baseUrl, stopPointPrefix, v.StopPoints[0].ShortName)
			destinationStop = fmt.Sprintf("%v%v/%v", baseUrl, stopPointPrefix, v.StopPoints[len(v.StopPoints)-1].ShortName)
		}
		converted.JourneyPatterns = append(converted.JourneyPatterns, RouteJourneyPattern{
			Url:             fmt.Sprintf("%v%v/%v", baseUrl, journeyPatternPrefix, v.Id),
			OriginStop:      originStop,
			DestinationStop: destinationStop,
			Name:            name,
		})
	}

	return converted
}

type Route struct {
	Projection      string                `json:"geographicCoordinateProjection"`
	Url             string                `json:"url"`
	Name            string                `json:"name"`
	LineUrl         string                `json:"lineUrl"`
	JourneyPatterns []RouteJourneyPattern `json:"journeyPatterns"`
	Journeys        []RouteJourney        `json:"journeys"`
}

type RouteJourneyPattern struct {
	Url             string `json:"url"`
	OriginStop      string `json:"originStop"`
	DestinationStop string `json:"destinationStop"`
	Name            string `json:"name"`
}

type RouteJourney struct {
	Url               string             `json:"url"`
	JourneyPatternUrl string             `json:"journeyPatternUrl"`
	DepartureTime     string             `json:"departureTime"`
	ArrivalTime       string             `json:"arrivalTime"`
	DayTypes          []string           `json:"dayTypes"`
	DayTypeExceptions []DayTypeException `json:"dayTypeExceptions"`
}
