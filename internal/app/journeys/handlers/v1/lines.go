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
		modelLines := service.SearchLines(getQueryParameters(req))

		var lines []Line
		for _, ml := range modelLines {
			lines = append(lines, convertLine(ml, baseUrl))
		}

		lex, err := removeExcludedFields(lines, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(lines)), rw)
		}

		sendJson(newSuccessResponse(lex), rw)
	}
}

func HandleGetOneLine(service service.DataService, baseUrl string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ml, err := service.GetOneLineById(mux.Vars(req)["name"])
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(make([]Line, 0))), rw)
			return
		}

		lines := []Line{convertLine(ml, baseUrl)}

		lex, err := removeExcludedFields(lines, getExcludeFieldsQueryParameter(req))
		if err != nil {
			sendJson(newSuccessResponse(arrayToAnyArray(lines)), rw)
		}

		sendJson(newSuccessResponse(lex), rw)
	}
}

func convertLine(line *model.Line, baseUrl string) Line {
	return Line{
		Url:         fmt.Sprintf("%v%v/%v", baseUrl, linePrefix, line.Name),
		Name:        line.Name,
		Description: line.Description,
	}
}

type Line struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
