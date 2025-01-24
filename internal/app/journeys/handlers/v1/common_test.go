//go:build journeys_common_tests || journeys_tests || all_tests

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

type successResponse[T APIEntity] struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []T            `json:"body"`
}

type handlerConfig struct {
	handler http.HandlerFunc
	url     string
}
type routerTestCase[T APIEntity] struct {
	target           string
	expectedEntities []T
	errorExpected    bool
	handlerConfig    handlerConfig
}

func newJourneysTestDataService() service.DataService {
	return service.DataService{DataStore: repository.NewJourneysDataStore("testdata/tre/gtfs", true)}
}

func runRouterTestCases[E APIEntity](t *testing.T, testCases []routerTestCase[E]) {
	for _, tc := range testCases {
		router := mux.NewRouter()
		router.HandleFunc(tc.handlerConfig.url, tc.handlerConfig.handler)

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

		bb := rec.Body.Bytes()
		var gotResponse successResponse[E]
		err := json.Unmarshal(bb, &gotResponse)
		if err != nil {
			t.Error(err)
			t.Log(string(bb))
		}

		expectedResponse := successResponse[E]{
			Status: "success",
			Data: apiSuccessData{
				Headers: apiHeaders{
					Paging: apiHeadersPaging{
						StartIndex: 0,
						PageSize:   len(tc.expectedEntities),
						MoreData:   false,
					},
				},
			},
			Body: tc.expectedEntities,
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

//var mockFilterObject, originalFilterObject func(r interface{}, properties string) (interface{}, error)
//var fakeMarshal, originalMarshal func(_ interface{}) ([]byte, error)
//
//func init() {
//	_ = os.Setenv("JOURNEYS_GTFS_PATH", "testdata/tre/gtfs")
//
//	mockFilterObject = func(r interface{}, properties string) (interface{}, error) {
//		return nil, errors.New("foo")
//	}
//
//	fakeMarshal = func(_ interface{}) ([]byte, error) {
//		return []byte{}, errors.New("marshalling failed")
//	}
//}
//
////func TestCreateSuccessResponseFail(t *testing.T) {
////	originalFilterObject = filterObject
////	filterObject = mockFilterObject
////	_, err := createSuccessResponse(make([]interface{}, 1), 0, 0, false, "")
////	if err == nil {
////		t.Error("Expected to see an error")
////	}
////
////	filterObject = originalFilterObject
////}
////
////func TestSendResponseFail(t *testing.T) {
////	w := &responseWriterMock{}
////	sendResponse(nil, errors.New("foobar"), nil, w)
////	if string(w.BodyData) != `{"status":"fail","data":{"message":"foobar"}}` {
////		t.Error("expected an error message")
////	}
////
////	w.Reset()
////
////	originalFilterObject = filterObject
////	filterObject = mockFilterObject
////	sendResponse(make([]interface{}, 1), nil, nil, w)
////	if string(w.BodyData) != `{"status":"fail","data":{"message":"foo"}}` {
////		t.Error("expected an error message")
////	}
////	filterObject = originalFilterObject
////
////	w.Reset()
////
////	originalMarshal = jsonMarshal
////	jsonMarshal = fakeMarshal
////	sendResponse(nil, nil, nil, w)
////	if w.StatusCode != 500 {
////		t.Error("expected status code: 500")
////	}
////	jsonMarshal = originalMarshal
////
////	w.Reset()
////	w.ErrorOnWrite = true
////
////	sendResponse(nil, nil, nil, w)
////	if w.StatusCode != 500 {
////		t.Error("expected status code: 500")
////	}
////	w.Reset()
////
////}
//
////func initializeTest(t *testing.T) (*mux.Router, *httptest.ResponseRecorder) {
////	dataService := service.DataService{DataStore: repository.NewJourneysDataStore("testdata/tre/gtfs", true)}
////	r := router.New(dataService, "http://localhost:8080")
////	w := httptest.NewRecorder()
////
////	return r, w
////}
//
//func serveHttp(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) {
//	r.ServeHTTP(w, httptest.NewRequest("GET", target, nil))
//
//	if w.Code != http.StatusOK {
//		t.Error("Did not get expected HTTP status code, got", w.Code)
//	}
//}
//
//func lineUrl(name string) string {
//	return fmt.Sprintf("%v/lines/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
//}
//
//func journeyPatternUrl(name string) string {
//	return fmt.Sprintf("%v/journey-patterns/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
//}
//
//func journeyUrl(name string) string {
//	return fmt.Sprintf("%v/journeys/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
//}
//
//func routeUrl(name string) string {
//	return fmt.Sprintf("%v/routes/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
//}
//
//func journeyActivityUrl(name string) string {
//	return fmt.Sprintf("%v/vehicle-activity/%v", os.Getenv("JOURNEYS_VA_BASE_URL"), name)
//}
//
//type responseWriterMock struct {
//	BodyData     []byte
//	StatusCode   int
//	ErrorOnWrite bool
//}
//
//func (rw *responseWriterMock) Header() http.Header {
//	return http.Header{}
//}
//
//func (rw *responseWriterMock) Write(bytes []byte) (int, error) {
//	if rw.ErrorOnWrite {
//		return 0, errors.New("write error")
//	}
//
//	rw.BodyData = bytes
//	return len(bytes), nil
//}
//
//func (rw *responseWriterMock) WriteHeader(statusCode int) {
//	rw.StatusCode = statusCode
//}
//
//func (rw *responseWriterMock) Reset() {
//	rw.StatusCode = 0
//	rw.BodyData = nil
//	rw.ErrorOnWrite = false
//}
