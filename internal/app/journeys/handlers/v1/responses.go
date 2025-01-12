package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIEntity interface {
	Line | Journey | JourneyPattern | Route | StopPoint | Municipality
}

func sendJson(data any, w http.ResponseWriter) {
	response, err := json.Marshal(data)
	if err != nil {
		sendError(err.Error(), w)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		sendError(err.Error(), w)
		return
	}
}

func sendError(errMsg string, w http.ResponseWriter) {
	afr := apiFailResponse{
		Status: "fail",
		Data: apiFailData{
			Message: fmt.Sprintf("%s", errMsg),
		},
	}
	response, jErr := json.Marshal(afr)
	if jErr != nil {
		log.Println(jErr)
		http.Error(w, jErr.Error(), http.StatusInternalServerError)
		return
	}

	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newSuccessResponse(body []any) apiSuccessResponse {
	return apiSuccessResponse{
		Status: "success",
		Data: apiSuccessData{Headers: apiHeaders{Paging: apiHeadersPaging{
			StartIndex: 0,
			PageSize:   len(body),
			MoreData:   false,
		}}},
		Body: body,
	}
}

type apiSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []any          `json:"body"`
}

type apiSuccessData struct {
	Headers apiHeaders `json:"headers"`
}

type apiHeaders struct {
	Paging apiHeadersPaging `json:"paging"`
}

type apiHeadersPaging struct {
	StartIndex int  `json:"startIndex"`
	PageSize   int  `json:"pageSize"`
	MoreData   bool `json:"moreData"`
}

type apiFailResponse struct {
	Status string      `json:"status"`
	Data   apiFailData `json:"data"`
}
type apiFailData struct {
	Message string `json:"message"`
}
