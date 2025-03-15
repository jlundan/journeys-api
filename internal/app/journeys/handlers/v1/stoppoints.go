package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllStopPoints(service *service.JourneysDataService, baseUrl string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelStopPoints := service.StopPoints.Search(getQueryParameters(req))

		var stopPoints []StopPoint
		for _, msp := range modelStopPoints {
			stopPoints = append(stopPoints, convertStopPoint(msp, baseUrl))
		}

		sendSuccessResponse(stopPoints, getExcludeFieldsQueryParameter(req), rw)
	}
}

func HandleGetOneStopPoint(service *service.JourneysDataService, baseUrl string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		msp, err := service.StopPoints.GetOneById(mux.Vars(req)["name"])
		if err != nil {
			sendSuccessResponse([]StopPoint{}, getExcludeFieldsQueryParameter(req), rw)
			return
		}

		stopPoints := []StopPoint{convertStopPoint(msp, baseUrl)}
		sendSuccessResponse(stopPoints, getExcludeFieldsQueryParameter(req), rw)
	}
}

func HandleGetJourneysForStopPoint(service *service.JourneysDataService, baseUrl string, vehicleActivityBaseUrl string, excludeInactive bool) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		stopPointId := mux.Vars(req)["name"]

		searchParams := map[string]string{}
		searchParams["stopPointId"] = stopPointId
		modelJourneys := service.Journeys.Search(searchParams, excludeInactive)

		var stopPointJourneys []StopPointJourney
		for _, mj := range modelJourneys {
			stopPointJourneys = append(stopPointJourneys, convertStopPointJourney(stopPointId, mj, baseUrl, vehicleActivityBaseUrl))
		}

		sendSuccessResponse(stopPointJourneys, getExcludeFieldsQueryParameter(req), rw)
	}
}

func convertStopPoint(stopPoint *model.StopPoint, baseUrl string) StopPoint {
	return StopPoint{
		Url:          fmt.Sprintf("%v%v/%v", baseUrl, stopPointPrefix, stopPoint.ShortName),
		ShortName:    stopPoint.ShortName,
		Name:         stopPoint.Name,
		Location:     fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude),
		TariffZone:   stopPoint.TariffZone,
		Municipality: convertStopPointMunicipality(stopPoint.Municipality, baseUrl),
	}
}

func convertStopPointMunicipality(municipality *model.Municipality, baseUrl string) StopPointMunicipality {
	return StopPointMunicipality{
		Url:       fmt.Sprintf("%v%v/%v", baseUrl, municipalitiesPrefix, municipality.PublicCode),
		ShortName: municipality.PublicCode,
		Name:      municipality.Name,
	}
}

func convertStopPointJourney(stopPointId string, j *model.Journey, baseUrl string, vehicleActivityBaseUrl string) StopPointJourney {
	var arrivalTime, departureTime string
	for _, c := range j.Calls {
		if c.StopPoint != nil && c.StopPoint.ShortName == stopPointId {
			arrivalTime = c.ArrivalTime
			departureTime = c.DepartureTime
		}
	}

	dayTypeExceptions := makeStopJourneyDayTypeExceptions(j)

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

	var gtfsInfo StopPointJourneyGtfsInfo

	if j.GtfsInfo != nil {
		gtfsInfo = StopPointJourneyGtfsInfo{TripId: j.GtfsInfo.TripId}
	}

	return StopPointJourney{
		JourneyUrl:           fmt.Sprintf("%v%v/%v", baseUrl, journeysPrefix, j.Id),
		StopPointUrl:         fmt.Sprintf("%v%v/%v", baseUrl, stopPointPrefix, stopPointId),
		ActivityUrl:          fmt.Sprintf("%v%v?journeyRef=%v", vehicleActivityBaseUrl, "/vehicle-activity", j.ActivityId),
		HeadSign:             j.HeadSign,
		Direction:            j.Direction,
		WheelchairAccessible: j.WheelchairAccessible,
		GtfsInfo:             gtfsInfo,
		JourneyPatternUrl:    fmt.Sprintf("%v%v/%v", baseUrl, journeyPatternPrefix, journeyPatternId),
		LineId:               lineId,
		LineUrl:              fmt.Sprintf("%v%v/%v", baseUrl, linePrefix, lineId),
		RouteUrl:             fmt.Sprintf("%v%v/%v", baseUrl, routePrefix, routeId),
		DayTypes:             j.DayTypes,
		DayTypeExceptions:    dayTypeExceptions,
		DepartureTime:        departureTime,
		ArrivalTime:          arrivalTime,
		ActiveFrom:           j.ValidFrom,
		ActiveTo:             j.ValidTo,
	}
}

type StopPoint struct {
	Url          string                `json:"url"`
	ShortName    string                `json:"shortName"`
	Name         string                `json:"name"`
	Location     string                `json:"location"`
	TariffZone   string                `json:"tariffZone"`
	Municipality StopPointMunicipality `json:"municipality"`
}

type StopPointMunicipality struct {
	Url       string `json:"url"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

type StopPointJourney struct {
	JourneyUrl           string                      `json:"journeyUrl"`
	StopPointUrl         string                      `json:"stopPointUrl"`
	ActivityUrl          string                      `json:"activityUrl"`
	LineUrl              string                      `json:"lineUrl"`
	RouteUrl             string                      `json:"routeUrl"`
	JourneyPatternUrl    string                      `json:"journeyPatternUrl"`
	LineId               string                      `json:"lineId"`
	DepartureTime        string                      `json:"departureTime"`
	ArrivalTime          string                      `json:"arrivalTime"`
	HeadSign             string                      `json:"headSign"`
	Direction            string                      `json:"directionId"`
	WheelchairAccessible bool                        `json:"wheelchairAccessible"`
	GtfsInfo             StopPointJourneyGtfsInfo    `json:"gtfs"`
	DayTypes             []string                    `json:"dayTypes"`
	DayTypeExceptions    []StopPointDayTypeException `json:"dayTypeExceptions"`
	ActiveFrom           string                      `json:"activeFrom"`
	ActiveTo             string                      `json:"activeTo"`
}

type StopPointJourneyGtfsInfo struct {
	TripId string `json:"tripId"`
}

type StopPointDayTypeException struct {
	From string `json:"from"`
	To   string `json:"to"`
	Runs string `json:"runs"`
}

func makeStopJourneyDayTypeExceptions(journey *model.Journey) []StopPointDayTypeException {
	dayTypeExceptions := make([]StopPointDayTypeException, 0)
	for _, dte := range journey.DayTypeExceptions {
		var runs string
		if dte.Runs {
			runs = "yes"
		} else {
			runs = "no"
		}

		dayTypeExceptions = append(dayTypeExceptions, StopPointDayTypeException{
			From: dte.From,
			To:   dte.To,
			Runs: runs,
		})
	}
	return dayTypeExceptions
}
