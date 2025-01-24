package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllMunicipalities(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelMunicipalities := service.SearchMunicipalities(getQueryParameters(req))

		var municipalities []Municipality
		for _, mm := range modelMunicipalities {
			municipalities = append(municipalities, convertMunicipality(mm, baseUrl))
		}

		lem, err := removeExcludedFields(municipalities, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(municipalities)), rw)
		}

		sendJson(newSuccessResponse(lem), rw)
	}
}

func HandleGetOneMunicipality(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		mm, err := service.GetOneMunicipalityById(mux.Vars(req)["name"])
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(make([]Municipality, 0))), rw)
			return
		}

		municipalities := []Municipality{convertMunicipality(mm, baseUrl)}

		mex, err := removeExcludedFields(municipalities, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(municipalities)), rw)
		}

		sendJson(newSuccessResponse(mex), rw)
	}
}

func convertMunicipality(municipality *model.Municipality, baseUrl string) Municipality {
	return Municipality{
		Url:       fmt.Sprintf("%v%v/%v", baseUrl, municipalitiesPrefix, municipality.PublicCode),
		ShortName: municipality.PublicCode,
		Name:      municipality.Name,
	}
}

type Municipality struct {
	Url       string `json:"url"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}
