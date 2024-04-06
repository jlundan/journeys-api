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
	var direction string
	var lineUrl string

	if l := len(route.Journeys); l > 0 {
		direction = route.Journeys[0].Direction
		lineUrl = convertJourney(route.Journeys[0]).LineUrl
	}

	converted := Route{
		Projection: route.GeoProjection,
		Direction:  direction,
		Links: map[string]string{
			"self":             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, route.Id),
			"line":             lineUrl,
			"journey-patterns": fmt.Sprintf("%v%v?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, route.Id),
			"journeys":         fmt.Sprintf("%v%v?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, route.Id),
		},
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
	Projection string            `json:"geographicCoordinateProjection"`
	Direction  string            `json:"direction"`
	Links      map[string]string `json:"links"`
}
