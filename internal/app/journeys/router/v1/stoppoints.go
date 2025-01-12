package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllStopPoints(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelStopPoints := service.SearchStopPoints(getQueryParameters(req))

		var stopPoints []StopPoint
		for _, msp := range modelStopPoints {
			stopPoints = append(stopPoints, convertStopPoint(msp, baseUrl))
		}

		spEx, err := removeExcludedFields(stopPoints, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(stopPoints)), rw)
		}

		sendJson(newSuccessResponse(spEx), rw)
	}
}

func HandleGetOneStopPoint(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		msp, err := service.GetOneStopPointById(mux.Vars(req)["name"])
		if err != nil {
			sendError("no such element", rw)
			return
		}

		stopPoints := []StopPoint{convertStopPoint(msp, baseUrl)}

		spEx, err := removeExcludedFields(stopPoints, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(stopPoints)), rw)
		}

		sendJson(newSuccessResponse(spEx), rw)
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
