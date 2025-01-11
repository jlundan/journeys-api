package v1

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"net/http"
)

func HandleGetAllLines(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		qp := getQueryParameters(req)
		lines := convertLines(service.SearchLines(qp), baseUrl)

		lex, err := removeExcludedFields(lines, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(lines)), rw)
		}

		sendJson(newSuccessResponse(lex), rw)
	}
}

func HandleGetOneLine(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		line := service.GetOneLineByName(mux.Vars(req)["name"])
		if line == nil {
			sendError("no such element", rw)
			return
		}

		lines := convertLines([]*model.Line{line}, baseUrl)

		lex, err := removeExcludedFields(lines, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(lines)), rw)
		}

		sendJson(newSuccessResponse(lex), rw)
	}
}

func convertLines(lines []*model.Line, baseUrl string) []Line {
	if len(lines) == 0 {
		return []Line{}
	}

	var converted []Line
	for _, line := range lines {
		converted = append(converted, Line{
			Url:         fmt.Sprintf("%v%v/%v", baseUrl, "/lines", line.Name),
			Name:        line.Name,
			Description: line.Description,
		})
	}
	return converted
}

type Line struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
