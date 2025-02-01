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
