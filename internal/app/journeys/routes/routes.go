package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/utils"
	"net/http"
	"os"
)

const routePrefix = "/routes"

func InjectRouteRoutes(r *mux.Router, context model.Context) {
	sr := r.PathPrefix(routePrefix).Subrouter()

	sr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllRoutes(w, r, context)
	}).Methods("GET")

	sr.HandleFunc(`/{name}`, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handleGetOneRoute(w, r, context, params["name"])
	}).Methods("GET")
}

func handleGetAllRoutes(w http.ResponseWriter, r *http.Request, context model.Context) {
	responseItems := make([]interface{}, 0)

	for _, route := range context.Routes().GetAll() {
		if routeMatchesConditions(route, getDefaultConditions(r)) {
			responseItems = append(responseItems, convertRoute(route))
		}
	}

	sendResponse(responseItems, nil, r, w)
}

func handleGetOneRoute(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
	route, err := context.Routes().GetOne(name)

	if err != nil {
		sendResponse(nil, err, r, w)
	} else {
		sendResponse([]interface{}{convertRoute(route)}, nil, r, w)
	}
}

func convertRoute(route *model.Route) Route {
	var name string

	if len(route.JourneyPatterns) > 0 && len(route.JourneyPatterns[0].StopPoints) > 0 {
		var firstStopPoint = route.JourneyPatterns[0].StopPoints[0]
		var lastStopPoint = route.JourneyPatterns[0].StopPoints[len(route.JourneyPatterns[0].StopPoints)-1]
		name = fmt.Sprintf("%v - %v", firstStopPoint.Name, lastStopPoint.Name)
	} else {
		name = ""
	}

	converted := Route{
		Projection: route.GeoProjection,
		Url:        fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, route.Id),
		LineUrl:    fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), linePrefix, route.Line.Name),
		Name:       name,
		//Links: map[string]string{
		//	"self":             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, route.Id),
		//	"line":             lineUrl,
		//	"journey-patterns": fmt.Sprintf("%v%v?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, route.Id),
		//	"journeys":         fmt.Sprintf("%v%v?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, route.Id),
		//},
	}

	for _, v := range route.Journeys {
		converted.Journeys = append(converted.Journeys, RouteJourney{
			Url:               fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, v.Id),
			JourneyPatternUrl: fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, v.JourneyPattern.Id),
			DepartureTime:     v.DepartureTime,
			ArrivalTime:       v.ArrivalTime,
			DayTypes:          v.DayTypes,
			DayTypeExceptions: makeDayTypeExceptions(v),
		})
	}

	for _, v := range route.JourneyPatterns {
		var originStop, destinationStop string

		if len(v.StopPoints) > 0 {
			originStop = fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, v.StopPoints[0].ShortName)
			destinationStop = fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, v.StopPoints[len(v.StopPoints)-1].ShortName)
		}
		converted.JourneyPatterns = append(converted.JourneyPatterns, RouteJourneyPattern{
			Url:             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, v.Id),
			OriginStop:      originStop,
			DestinationStop: destinationStop,
			Name:            name,
		})
	}

	return converted
}

func routeMatchesConditions(route *model.Route, conditions map[string]string) bool {
	if route == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(route.Name, v) {
				return false
			}
		case "lineId":
			if route.Line == nil || route.Line.Name != v {
				return false
			}
		}
	}

	return true
}

type Route struct {
	Projection string `json:"geographicCoordinateProjection"`
	//	Direction       string   `json:"direction"`
	Url             string                `json:"url"`
	Name            string                `json:"name"`
	LineUrl         string                `json:"lineUrl"`
	JourneyPatterns []RouteJourneyPattern `json:"journeyPatterns"`
	Journeys        []RouteJourney        `json:"journeys"`
	//Links      	map[string]string 	`json:"links"`
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
