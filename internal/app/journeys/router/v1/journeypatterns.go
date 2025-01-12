package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
	"os"
)

func HandleGetAllJourneyPatterns(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelJourneyPatterns := service.SearchJourneyPatterns(getQueryParameters(req))

		var journeyPatterns []JourneyPattern
		for _, mjp := range modelJourneyPatterns {
			journeyPatterns = append(journeyPatterns, convertJourneyPattern(mjp, baseUrl))
		}

		jpe, err := removeExcludedFields(journeyPatterns, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(journeyPatterns)), rw)
		}

		sendJson(newSuccessResponse(jpe), rw)
	}
}

func HandleGetOneJourneyPattern(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		mj, err := service.GetOneJourneyPatternById(mux.Vars(req)["name"])
		if err != nil {
			sendError("no such element", rw)
			return
		}

		journeyPatterns := []JourneyPattern{convertJourneyPattern(mj, baseUrl)}

		jex, err := removeExcludedFields(journeyPatterns, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(journeyPatterns)), rw)
		}

		sendJson(newSuccessResponse(jex), rw)
	}
}

func convertJourneyPattern(jp *model.JourneyPattern, baseUrl string) JourneyPattern {
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
		converted.StopPoints = append(converted.StopPoints, convertJourneyPatternStopPoint(v))
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

func convertJourneyPatternStopPoint(stopPoint *model.StopPoint) JourneyPatternStopPoint {
	return JourneyPatternStopPoint{
		Url:          fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, stopPoint.ShortName),
		ShortName:    stopPoint.ShortName,
		Name:         stopPoint.Name,
		Location:     fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude),
		TariffZone:   stopPoint.TariffZone,
		Municipality: convertJourneyPatternMunicipality(stopPoint.Municipality),
	}
}

func convertJourneyPatternMunicipality(municipality *model.Municipality) JourneyPatternMunicipality {
	return JourneyPatternMunicipality{
		Url:       fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), municipalitiesPrefix, municipality.PublicCode),
		ShortName: municipality.PublicCode,
		Name:      municipality.Name,
	}
}

type JourneyPattern struct {
	Url             string                    `json:"url"`
	RouteUrl        string                    `json:"routeUrl"`
	LineUrl         string                    `json:"lineUrl"`
	OriginStop      string                    `json:"originStop"`
	DestinationStop string                    `json:"destinationStop"`
	Name            string                    `json:"name"`
	StopPoints      []JourneyPatternStopPoint `json:"stopPoints"`
	Journeys        []JourneyPatternJourney   `json:"journeys"`
	Direction       string                    `json:"direction"`
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

type JourneyPatternStopPoint struct {
	Url          string                     `json:"url"`
	ShortName    string                     `json:"shortName"`
	Name         string                     `json:"name"`
	Location     string                     `json:"location"`
	TariffZone   string                     `json:"tariffZone"`
	Municipality JourneyPatternMunicipality `json:"municipality"`
}

type JourneyPatternMunicipality struct {
	Url       string `json:"url"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}
