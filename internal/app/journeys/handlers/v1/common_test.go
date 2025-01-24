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
