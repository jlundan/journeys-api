//go:build journeys_common_tests || journeys_tests || all_tests

package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"github.com/jlundan/journeys-api/internal/testutil"
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

func newJourneysTestDataService(t *testing.T) *service.JourneysDataService {
	repo, errs := repository.NewJourneysRepository("testdata/tre/gtfs", true)
	if len(errs) > 0 {
		t.Error(errs)
	}
	return service.NewJourneysDataService(repo)
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

		var diffs []testutil.FieldDiff
		initialTag := fmt.Sprintf("%v:Response", tc.target)
		err = testutil.CompareVariables(expectedResponse, gotResponse, initialTag, &diffs, false)
		if err != nil {
			t.Error(err)
			break
		}

		if len(diffs) > 0 {
			testutil.PrintFieldDiffs(t, diffs)
			break
		}
	}
}
