//go:build journeys_municipalities_tests || journeys_tests || all_tests

package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSendJson(t *testing.T) {
	// Test successful JSON response
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	sendJson(data, rr)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("sendJson() status code = %v, want %v", status, http.StatusOK)
	}
	expected, _ := json.Marshal(data)
	if rr.Body.String() != string(expected) {
		t.Errorf("sendJson() body = %v, want %v", rr.Body.String(), string(expected))
	}

	// Test JSON marshalling error
	rr = httptest.NewRecorder()
	sendJson(make(chan int), rr)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("sendJson() status code = %v, want %v", status, http.StatusInternalServerError)
	}

}

type errorReturningHttpWriter struct {
	Code int
}

func (e *errorReturningHttpWriter) Header() http.Header {
	return http.Header{}
}
func (e *errorReturningHttpWriter) WriteHeader(statusCode int) {
	e.Code = statusCode
}
func (e *errorReturningHttpWriter) Write(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, errors.New("test error")
	}
	return 0, nil
}

func TestSendResponse(t *testing.T) {
	erw := errorReturningHttpWriter{}
	sendResponse([]byte{}, &erw)
	if status := erw.Code; status != http.StatusInternalServerError {
		t.Errorf("sendResponse() status code = %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestNewSuccessResponse(t *testing.T) {
	// Test successful response
	body := []map[string]interface{}{{"key": "value"}}
	expected := apiSuccessResponse{"success", apiSuccessData{apiHeaders{apiHeadersPaging{
		StartIndex: 0,
		PageSize:   1,
		MoreData:   false,
	}}}, body}
	if response := newSuccessResponse(body); !reflect.DeepEqual(response, expected) {
		t.Errorf("newSuccessResponse() = %v, want %v", response, expected)
	}

	// Test that nil bodies are returned as empty arrays
	expected = apiSuccessResponse{"success", apiSuccessData{apiHeaders{apiHeadersPaging{
		StartIndex: 0,
		PageSize:   0,
		MoreData:   false,
	}}}, []map[string]interface{}{}}
	if response := newSuccessResponse(nil); !reflect.DeepEqual(response, expected) {
		t.Errorf("newSuccessResponse() = %v, want %v", response, expected)
	}
}
