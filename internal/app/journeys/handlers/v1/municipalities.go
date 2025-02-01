package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllMunicipalities(service *service.JourneysDataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		modelMunicipalities := service.Municipalities.Search(getQueryParameters(req))

		var municipalities []Municipality
		for _, mm := range modelMunicipalities {
			municipalities = append(municipalities, convertMunicipality(mm, baseUrl))
		}

		sendSuccessResponse(municipalities, getExcludeFieldsQueryParameter(req), rw)
	}
}

func HandleGetOneMunicipality(service *service.JourneysDataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		mm, err := service.Municipalities.GetOneById(mux.Vars(req)["name"])
		if err != nil {
			sendSuccessResponse([]Municipality{}, getExcludeFieldsQueryParameter(req), rw)
			return
		}

		municipalities := []Municipality{convertMunicipality(mm, baseUrl)}
		sendSuccessResponse(municipalities, getExcludeFieldsQueryParameter(req), rw)
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
