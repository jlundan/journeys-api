//go:build journeys_stops_tests || journeys_tests || all_tests

package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStopPointRoutes(t *testing.T) {
	dataService := service.DataService{DataStore: repository.NewJourneysDataStore("testdata/tre/gtfs", true)}

	oneHandler := HandleGetOneStopPoint(dataService, "")
	oneHandlerUrl := "/v1/stop-points/{name}"
	allHandler := HandleGetAllStopPoints(dataService, "")
	allHandlerUrl := "/v1/stop-points"

	var stopPointMap = getStopPointMap()
	testCases := []struct {
		target        string
		items         []StopPoint
		errorExpected bool
		handler       http.HandlerFunc
		handlerUrl    string
	}{
		{"/v1/stop-points",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"], stopPointMap["4600"],
				stopPointMap["7015"], stopPointMap["7017"], stopPointMap["8149"], stopPointMap["8171"],
			}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?tariffZone=B",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"], stopPointMap["4600"],
				stopPointMap["7015"], stopPointMap["7017"], stopPointMap["8171"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?municipalityName=Tampere",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?municipalityShortName=837",
			[]StopPoint{stopPointMap["3607"], stopPointMap["3615"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?location=61.4445,23.87235",
			[]StopPoint{stopPointMap["3615"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?location=61,23:62,23.8",
			[]StopPoint{stopPointMap["7015"], stopPointMap["7017"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?name=Pirkkala",
			[]StopPoint{stopPointMap["7015"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?name=Sudenkorennontie&shortName=81",
			[]StopPoint{stopPointMap["8149"]}, false, allHandler, allHandlerUrl,
		},
		{"/v1/stop-points?shortName=4600&exclude-fields=name,shortName,municipality.shortName", []StopPoint{
			{
				stopPointUrl("4600"),
				"",
				"",
				"61.47561,23.97756",
				"B",
				StopPointMunicipality{
					Url:  stopPointMunicipalityUrl("211"),
					Name: "Kangasala",
				},
			},
		}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points/4600", []StopPoint{stopPointMap["4600"]}, false, oneHandler, oneHandlerUrl},
		{"/v1/stop-points/foobar", []StopPoint{}, false, oneHandler, oneHandlerUrl},
		{"/v1/stop-points?location=6123:62,23.8", []StopPoint{}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?location=61,23:6223.8", []StopPoint{}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?location=A1,23:62,23.8", []StopPoint{}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?location=61,23:62,23.X", []StopPoint{}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?location=61,2X:62,23.8", []StopPoint{}, false, allHandler, allHandlerUrl},
		{"/v1/stop-points?location=61,23:6X,23.8", []StopPoint{}, false, allHandler, allHandlerUrl},
	}

	for _, tc := range testCases {
		router := mux.NewRouter()
		router.HandleFunc(tc.handlerUrl, tc.handler)

		req := httptest.NewRequest("GET", tc.target, nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if tc.errorExpected {
			var response apiFailResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				t.Error(err)
			}
			continue
		}

		var gotResponse stopPointSuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &gotResponse)
		if err != nil {
			t.Error(err)
		}

		expectedResponse := stopPointSuccessResponse{
			Status: "success",
			Data: apiSuccessData{
				Headers: apiHeaders{
					Paging: apiHeadersPaging{
						StartIndex: 0,
						PageSize:   len(tc.items),
						MoreData:   false,
					},
				},
			},
			Body: tc.items,
		}

		var diffs []FieldDiff
		initialTag := fmt.Sprintf("%v:Response", tc.target)
		err = compareVariables(expectedResponse, gotResponse, initialTag, &diffs, false)
		if err != nil {
			t.Error(err)
			break
		}

		if len(diffs) > 0 {
			printFieldDiffs(t, diffs)
			break
		}
	}
}

func getStopPointMap() map[string]StopPoint {
	result := make(map[string]StopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality StopPointMunicipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getStopPointMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getStopPointMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getStopPointMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getStopPointMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getStopPointMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getStopPointMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getStopPointMunicipalityMap()["837"]},
	}

	for _, tc := range stopPoints {
		result[tc.id] = StopPoint{
			Url:          stopPointUrl(tc.id),
			ShortName:    tc.id,
			Name:         tc.name,
			Location:     tc.location,
			TariffZone:   tc.tariffZone,
			Municipality: tc.municipality,
		}
	}

	return result
}

func getStopPointMunicipalityMap() map[string]StopPointMunicipality {
	municipalities := make(map[string]StopPointMunicipality)
	municipalities["211"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("211"),
		ShortName: "211",
		Name:      "Kangasala",
	}

	municipalities["604"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("604"),
		ShortName: "604",
		Name:      "Pirkkala",
	}

	municipalities["837"] = StopPointMunicipality{
		Url:       stopPointMunicipalityUrl("837"),
		ShortName: "837",
		Name:      "Tampere",
	}
	return municipalities
}

type stopPointSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []StopPoint    `json:"body"`
}

func stopPointUrl(name string) string {
	return fmt.Sprintf("%v/stop-points/%v", "", name)
}
func stopPointMunicipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", "", name)
}
