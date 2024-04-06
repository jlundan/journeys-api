package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/utils"
	"net/http"
	"os"
)

const municipalitiesPrefix = "/municipalities"

func InjectMunicipalityRoutes(r *mux.Router, context model.Context) {
	sr := r.PathPrefix(municipalitiesPrefix).Subrouter()

	sr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllMunicipalities(w, r, context)
	}).Methods("GET")

	sr.HandleFunc(`/{name}`, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handleGetOneMunicipality(w, r, context, params["name"])
	}).Methods("GET")
}

func handleGetAllMunicipalities(w http.ResponseWriter, r *http.Request, context model.Context) {
	responseItems := make([]interface{}, 0)

	for _, municipality := range context.Municipalities().GetAll() {
		if municipalityMatchesConditions(municipality, getDefaultConditions(r)) {
			responseItems = append(responseItems, convertMunicipality(municipality))
		}
	}

	sendResponse(responseItems, nil, r, w)
}

func handleGetOneMunicipality(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
	municipality, err := context.Municipalities().GetOne(name)

	if err == nil {
		sendResponse([]interface{}{convertMunicipality(municipality)}, nil, r, w)
	} else {
		sendResponse(nil, err, r, w)
	}
}

func convertMunicipality(municipality *model.Municipality) Municipality {
	return Municipality{
		Url:       fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), municipalitiesPrefix, municipality.PublicCode),
		ShortName: municipality.PublicCode,
		Name:      municipality.Name,
	}
}

func municipalityMatchesConditions(municipality *model.Municipality, conditions map[string]string) bool {
	if municipality == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(municipality.Name, v) {
				return false
			}
		case "shortName":
			if !utils.StrContains(municipality.PublicCode, v) {
				return false
			}
		}
	}

	return true
}

type Municipality struct {
	Url       string `json:"url"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}
