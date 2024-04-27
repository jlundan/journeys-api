package routes

import (
	"encoding/json"
	"fmt"
	"github.com/jlundan/journeys-api/internal/pkg/utils"
	"log"
	"net/http"
)

var (
	filterObject func(r interface{}, properties string) (interface{}, error)
	jsonMarshal  func(_ interface{}) ([]byte, error)
)

func init() {
	filterObject = utils.FilterObject
	jsonMarshal = json.Marshal
}

func sendResponse(responseItems []interface{}, responseError error, r *http.Request, w http.ResponseWriter) {
	if responseError != nil && responseError.Error() != "no such element" {
		sendError(responseError, w)
		return
	}

	var resp interface{}
	var err error

	if responseItems == nil {
		resp, err = createSuccessResponse(make([]interface{}, 0), 0, 0, false, "")
	} else {
		var filter string
		if r != nil && r.URL != nil && r.URL.Query() != nil {
			filter = r.URL.Query().Get("exclude-fields")
		}
		resp, err = createSuccessResponse(responseItems, 0, uint16(len(responseItems)), false, filter)
	}

	if err != nil {
		sendError(err, w)
		return
	}

	err = sendJson(resp, w)
	if err != nil {
		sendError(err, w)
	}
}

func sendJson(data interface{}, w http.ResponseWriter) error {
	response, err := jsonMarshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func sendError(err error, w http.ResponseWriter) {
	afr := apiFailResponse{
		Status: "fail",
		Data: apiFailData{
			Message: fmt.Sprintf("%v", err),
		},
	}
	response, jErr := jsonMarshal(afr)
	if jErr != nil {
		log.Println(jErr)
		http.Error(w, jErr.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createSuccessResponse(body []interface{}, startIndex uint32, pageSize uint16, moreData bool, filter string) (apiSuccessResponse, error) {
	respArr := make([]interface{}, 0)
	for _, b := range body {
		b, err := filterObject(b, filter)
		if err != nil {
			return apiSuccessResponse{}, err
		}
		respArr = append(respArr, b)
	}

	return apiSuccessResponse{
		Status: "success",
		Data: apiSuccessData{Headers: apiHeaders{Paging: apiHeadersPaging{
			StartIndex: startIndex,
			PageSize:   pageSize,
			MoreData:   moreData,
		}}},
		Body: respArr,
	}, nil
}

func getDefaultConditions(r *http.Request) map[string]string {
	result := make(map[string]string)
	for k, v := range r.URL.Query() {
		if k != "exclude-fields" {
			result[k] = v[0]
		}
	}
	return result
}

type apiSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   interface{}    `json:"body"`
}

type apiSuccessData struct {
	Headers apiHeaders `json:"headers"`
}

type apiHeaders struct {
	Paging apiHeadersPaging `json:"paging"`
}

type apiHeadersPaging struct {
	StartIndex uint32 `json:"startIndex"`
	PageSize   uint16 `json:"pageSize"`
	MoreData   bool   `json:"moreData"`
}

type apiFailResponse struct {
	Status string      `json:"status"`
	Data   apiFailData `json:"data"`
}
type apiFailData struct {
	Message string `json:"message"`
}
