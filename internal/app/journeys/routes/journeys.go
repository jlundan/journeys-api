package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http"
	"os"
	"strings"
	"time"
)

const journeysPrefix = "/journeys"

func InjectJourneyRoutes(r *mux.Router, context model.Context) {
	sr := r.PathPrefix(journeysPrefix).Subrouter()

	sr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllJourneys(w, r, context)
	}).Methods("GET")

	sr.HandleFunc(`/{name}`, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handleGetOneJourney(w, r, context, params["name"])
	}).Methods("GET")
}

func handleGetAllJourneys(w http.ResponseWriter, r *http.Request, context model.Context) {
	responseItems := make([]interface{}, 0)

	for _, journey := range context.Journeys().GetAll() {
		if journeyMatchesConditions(journey, getDefaultConditions(r)) {
			responseItems = append(responseItems, convertJourney(journey))
		}
	}

	sendResponse(responseItems, nil, r, w)
}

func handleGetOneJourney(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
	jp, err := context.Journeys().GetOne(name)

	if err != nil && err.Error() == "no such element" {
		jp, err = context.Journeys().GetOneByActivityId(name)
	}

	if err != nil {
		sendResponse(nil, err, r, w)
	} else {
		sendResponse([]interface{}{convertJourney(jp)}, nil, r, w)
	}
}

func convertJourney(j *model.Journey) Journey {
	calls := make([]JourneyCall, 0)
	for _, c := range j.Calls {
		calls = append(calls, JourneyCall{
			DepartureTime: c.DepartureTime,
			ArrivalTime:   c.ArrivalTime,
			StopPoint:     convertStopPoint(c.StopPoint),
		})
	}

	dayTypeExceptions := makeDayTypeExceptions(j)

	var lineId, routeId, journeyPatternId string

	if j.Line != nil {
		lineId = j.Line.Name
	}
	if j.Route != nil {
		routeId = j.Route.Id
	}
	if j.JourneyPattern != nil {
		journeyPatternId = j.JourneyPattern.Id
	}

	var gtfsInfo JourneyGtfsInfo

	if j.GtfsInfo != nil {
		gtfsInfo = JourneyGtfsInfo{TripId: j.GtfsInfo.TripId}
	}

	return Journey{
		Url:                  fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, j.Id),
		ActivityUrl:          fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_VA_BASE_URL"), "/vehicle-activity", j.ActivityId),
		HeadSign:             j.HeadSign,
		Direction:            j.Direction,
		WheelchairAccessible: j.WheelchairAccessible,
		GtfsInfo:             gtfsInfo,
		JourneyPatternUrl:    fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, journeyPatternId),
		LineUrl:              fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), linePrefix, lineId),
		RouteUrl:             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, routeId),
		Calls:                calls,
		DayTypes:             j.DayTypes,
		DayTypeExceptions:    dayTypeExceptions,
		DepartureTime:        j.DepartureTime,
		ArrivalTime:          j.ArrivalTime,
	}
}

func convertSimpleJourney(j *model.Journey) SimpleJourney {
	dayTypeExceptions := makeDayTypeExceptions(j)

	var lineId, routeId, journeyPatternId string

	if j.Line != nil {
		lineId = j.Line.Name
	}
	if j.Route != nil {
		routeId = j.Route.Id
	}
	if j.JourneyPattern != nil {
		journeyPatternId = j.JourneyPattern.Id
	}

	var gtfsInfo JourneyGtfsInfo

	if j.GtfsInfo != nil {
		gtfsInfo = JourneyGtfsInfo{TripId: j.GtfsInfo.TripId}
	}

	return SimpleJourney{
		Url:                  fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeysPrefix, j.Id),
		ActivityUrl:          fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_VA_BASE_URL"), "/vehicle-activity", j.ActivityId),
		HeadSign:             j.HeadSign,
		Direction:            j.Direction,
		WheelchairAccessible: j.WheelchairAccessible,
		GtfsInfo:             gtfsInfo,
		JourneyPatternUrl:    fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), journeyPatternPrefix, journeyPatternId),
		LineUrl:              fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), linePrefix, lineId),
		RouteUrl:             fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), routePrefix, routeId),
		DayTypes:             j.DayTypes,
		DayTypeExceptions:    dayTypeExceptions,
		DepartureTime:        j.DepartureTime,
		ArrivalTime:          j.ArrivalTime,
	}
}

func journeyMatchesConditions(journey *model.Journey, conditions map[string]string) bool {
	now := time.Now()
	curDay := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

	if journey == nil || !(journey.ValidFrom <= curDay && journey.ValidTo >= curDay) {
		return false
	}

	if conditions == nil {
		return true
	}

	for k, v := range conditions {
		switch k {
		case "lineId":
			if journey.Line == nil || journey.Line.Name != v {
				return false
			}
		case "routeId":
			if journey.Route == nil || journey.Route.Id != v {
				return false
			}
		case "journeyPatternId":
			if journey.JourneyPattern == nil || journey.JourneyPattern.Id != v {
				return false
			}
		case "dayTypes":
			matched := false
			vDayTypes := strings.Split(v, ",")
			for _, dt := range journey.DayTypes {
				for _, vdt := range vDayTypes {
					if dt == vdt {
						matched = true
						break
					}
				}
			}
			if !matched {
				return false
			}
		case "departureTime":
			if journey.DepartureTime != v {
				return false
			}
		case "arrivalTime":
			if journey.ArrivalTime != v {
				return false
			}
		case "firstStopPointId":
			if len(journey.Calls) == 0 {
				return false
			}
			first := journey.Calls[0]
			if first.StopPoint == nil || first.StopPoint.ShortName != v {
				return false
			}
		case "lastStopPointId":
			if len(journey.Calls) == 0 {
				return false
			}
			last := journey.Calls[len(journey.Calls)-1]
			if last.StopPoint == nil || last.StopPoint.ShortName != v {
				return false
			}
		case "stopPointId":
			matched := false
			for _, c := range journey.Calls {
				if c.StopPoint != nil && c.StopPoint.ShortName == v {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		case "gtfsTripId":
			if journey.GtfsInfo == nil || journey.GtfsInfo.TripId != v {
				return false
			}
		}
	}

	return true
}

func makeDayTypeExceptions(journey *model.Journey) []DayTypeException {
	dayTypeExceptions := make([]DayTypeException, 0)
	for _, dte := range journey.DayTypeExceptions {
		var runs string
		if dte.Runs {
			runs = "yes"
		} else {
			runs = "no"
		}

		dayTypeExceptions = append(dayTypeExceptions, DayTypeException{
			From: dte.From,
			To:   dte.To,
			Runs: runs,
		})
	}
	return dayTypeExceptions
}

type Journey struct {
	Url                  string             `json:"url"`
	ActivityUrl          string             `json:"activityUrl"`
	LineUrl              string             `json:"lineUrl"`
	RouteUrl             string             `json:"routeUrl"`
	JourneyPatternUrl    string             `json:"journeyPatternUrl"`
	DepartureTime        string             `json:"departureTime"`
	ArrivalTime          string             `json:"arrivalTime"`
	HeadSign             string             `json:"headSign"`
	Direction            string             `json:"directionId"`
	WheelchairAccessible bool               `json:"wheelchairAccessible"`
	GtfsInfo             JourneyGtfsInfo    `json:"gtfs"`
	DayTypes             []string           `json:"dayTypes"`
	DayTypeExceptions    []DayTypeException `json:"dayTypeExceptions"`
	Calls                []JourneyCall      `json:"calls"`
}

type SimpleJourney struct {
	Url                  string             `json:"url"`
	ActivityUrl          string             `json:"activityUrl"`
	LineUrl              string             `json:"lineUrl"`
	RouteUrl             string             `json:"routeUrl"`
	JourneyPatternUrl    string             `json:"journeyPatternUrl"`
	DepartureTime        string             `json:"departureTime"`
	ArrivalTime          string             `json:"arrivalTime"`
	HeadSign             string             `json:"headSign"`
	Direction            string             `json:"directionId"`
	WheelchairAccessible bool               `json:"wheelchairAccessible"`
	GtfsInfo             JourneyGtfsInfo    `json:"gtfs"`
	DayTypes             []string           `json:"dayTypes"`
	DayTypeExceptions    []DayTypeException `json:"dayTypeExceptions"`
	//Calls                []JourneyCall      `json:"calls"`
}

type JourneyGtfsInfo struct {
	TripId string `json:"tripId"`
}

type DayTypeException struct {
	From string `json:"from"`
	To   string `json:"to"`
	Runs string `json:"runs"`
}

type JourneyCall struct {
	DepartureTime string    `json:"departureTime"`
	ArrivalTime   string    `json:"arrivalTime"`
	StopPoint     StopPoint `json:"stopPoint"`
}
