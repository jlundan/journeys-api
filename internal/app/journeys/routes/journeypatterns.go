package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/utils"
	"net/http"
	"os"
)

const journeyPatternPrefix = "/journey-patterns"

func InjectJourneyPatternRoutes(r *mux.Router, context model.Context) {
	sr := r.PathPrefix(journeyPatternPrefix).Subrouter()

	sr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllJourneyPatterns(w, r, context)
	}).Methods("GET")

	sr.HandleFunc(`/{name}`, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handleGetOneJourneyPattern(w, r, context, params["name"])
	}).Methods("GET")
}

func handleGetAllJourneyPatterns(w http.ResponseWriter, r *http.Request, context model.Context) {
	responseItems := make([]interface{}, 0)

	for _, journeyPattern := range context.JourneyPatterns().GetAll() {
		if journeyPatternMatchesConditions(journeyPattern, getDefaultConditions(r)) {
			responseItems = append(responseItems, convertJourneyPattern(journeyPattern))
		}
	}

	sendResponse(responseItems, nil, r, w)
}

func handleGetOneJourneyPattern(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
	jp, err := context.JourneyPatterns().GetOne(name)

	if err != nil {
		sendResponse(nil, err, r, w)
	} else {
		sendResponse([]interface{}{convertJourneyPattern(jp)}, nil, r, w)
	}
}

func convertJourneyPattern(jp *model.JourneyPattern) JourneyPattern {
	var direction string

	if len(jp.Route.Journeys) > 0 {
		direction = jp.Route.Journeys[0].Direction
	}

	var name string
	if len(jp.StopPoints) > 0 {
		name = fmt.Sprintf("%v - %v", jp.StopPoints[0].Name, jp.StopPoints[len(jp.StopPoints)-1].Name)
	} else {
		name = ""
	}

	converted := JourneyPattern{
		Url:             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, jp.Id),
		RouteUrl:        fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, jp.Route.Id),
		LineUrl:         fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), linePrefix, jp.Route.Line.Name),
		Name:            name,
		OriginStop:      fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, jp.StopPoints[0].ShortName),
		DestinationStop: fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, jp.StopPoints[len(jp.StopPoints)-1].ShortName),
		Direction:       direction,
	}

	for _, v := range jp.StopPoints {
		converted.StopPoints = append(converted.StopPoints, convertStopPoint(v))
	}

	for _, v := range jp.Journeys {
		converted.Journeys = append(converted.Journeys, JourneyPatternJourney{
			Url:               fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, v.Id),
			JourneyPatternUrl: fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, v.JourneyPattern.Id),
			DepartureTime:     v.DepartureTime,
			ArrivalTime:       v.ArrivalTime,
			HeadSign:          v.HeadSign,
			DayTypes:          v.DayTypes,
			DayTypeExceptions: makeDayTypeExceptions(v),
		})
	}
	return converted
}

func journeyPatternMatchesConditions(journeyPattern *model.JourneyPattern, conditions map[string]string) bool {
	if journeyPattern == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(journeyPattern.Name, v) {
				return false
			}
		case "lineId":
			if journeyPattern.Route == nil || journeyPattern.Route.Line == nil || journeyPattern.Route.Line.Name != v {
				return false
			}
		case "firstStopPointId":
			if len(journeyPattern.StopPoints) == 0 || journeyPattern.StopPoints[0].ShortName != v {
				return false
			}
		case "lastStopPointId":
			spLength := len(journeyPattern.StopPoints)
			if spLength == 0 || journeyPattern.StopPoints[spLength-1].ShortName != v {
				return false
			}
		case "stopPointId":
			found := false
			for _, sp := range journeyPattern.StopPoints {
				if sp.ShortName == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

type JourneyPattern struct {
	Url             string                  `json:"url"`
	RouteUrl        string                  `json:"routeUrl"`
	LineUrl         string                  `json:"lineUrl"`
	OriginStop      string                  `json:"originStop"`
	DestinationStop string                  `json:"destinationStop"`
	Name            string                  `json:"name"`
	StopPoints      []StopPoint             `json:"stopPoints"`
	Journeys        []JourneyPatternJourney `json:"journeys"`
	Direction       string                  `json:"direction"`
}

type JourneyPatternJourney struct {
	Url               string             `json:"url"`
	JourneyPatternUrl string             `json:"journeyPatternUrl"`
	DepartureTime     string             `json:"departureTime"`
	ArrivalTime       string             `json:"arrivalTime"`
	HeadSign          string             `json:"headSign"`
	DayTypes          []string           `json:"dayTypes"`
	DayTypeExceptions []DayTypeException `json:"dayTypeExceptions"`
}
