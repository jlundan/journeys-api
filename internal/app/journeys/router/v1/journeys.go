package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
	"os"
)

func HandleGetAllJourneys(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelJourneys := service.SearchJourneys(getQueryParameters(req))

		var journeys []Journey
		for _, mj := range modelJourneys {
			journeys = append(journeys, convertJourney(mj, baseUrl))
		}

		jex, err := removeExcludedFields(journeys, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(journeys)), rw)
		}

		sendJson(newSuccessResponse(jex), rw)
	}
}

func HandleGetOneJourney(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		mj, err := service.GetOneJourneyById(mux.Vars(req)["name"])
		if err != nil {
			sendError("no such element", rw)
			return
		}

		journeys := []Journey{convertJourney(mj, baseUrl)}

		jex, err := removeExcludedFields(journeys, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(journeys)), rw)
		}

		sendJson(newSuccessResponse(jex), rw)
	}
}

func convertJourney(j *model.Journey, baseUrl string) Journey {
	calls := make([]JourneyCall, 0)
	for _, c := range j.Calls {
		calls = append(calls, JourneyCall{
			DepartureTime: c.DepartureTime,
			ArrivalTime:   c.ArrivalTime,
			StopPoint:     convertJourneyStopPoint(c.StopPoint),
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

func convertJourneyStopPoint(stopPoint *model.StopPoint) JourneyStopPoint {
	return JourneyStopPoint{
		Url:          fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, stopPoint.ShortName),
		ShortName:    stopPoint.ShortName,
		Name:         stopPoint.Name,
		Location:     fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude),
		TariffZone:   stopPoint.TariffZone,
		Municipality: convertJourneyMunicipality(stopPoint.Municipality),
	}
}

func convertJourneyMunicipality(municipality *model.Municipality) JourneyMunicipality {
	return JourneyMunicipality{
		Url:       fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), municipalitiesPrefix, municipality.PublicCode),
		ShortName: municipality.PublicCode,
		Name:      municipality.Name,
	}
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

type JourneyGtfsInfo struct {
	TripId string `json:"tripId"`
}

type DayTypeException struct {
	From string `json:"from"`
	To   string `json:"to"`
	Runs string `json:"runs"`
}

type JourneyCall struct {
	DepartureTime string           `json:"departureTime"`
	ArrivalTime   string           `json:"arrivalTime"`
	StopPoint     JourneyStopPoint `json:"stopPoint"`
}

type JourneyStopPoint struct {
	Url          string              `json:"url"`
	ShortName    string              `json:"shortName"`
	Name         string              `json:"name"`
	Location     string              `json:"location"`
	TariffZone   string              `json:"tariffZone"`
	Municipality JourneyMunicipality `json:"municipality"`
}

type JourneyMunicipality struct {
	Url       string `json:"url"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}
