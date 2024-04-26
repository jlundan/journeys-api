//go:build journeys_lines_tests || journeys_tests || all_tests

package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"testing"
)

func TestLineMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Line
		conditions map[string]string
		expected   bool
	}{
		{
			nil,
			nil,
			false,
		},
		{
			&model.Line{Name: "1", Description: "Vatiala - Pirkkala"},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		matches := lineMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestGetLines(t *testing.T) {
	testCases := []struct {
		url   string
		items []Line
	}{
		{"/lines", []Line{
			{lineUrl("-1"), "-1", "Foobar"},
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
			{lineUrl("3A"), "3A", "Etelä-Hervanta - Lentävänniemi"},
		}},
		{"/lines?name=1", []Line{
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
		}},
		{"/lines?name=1A&description=lento", []Line{
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}},
		{"/lines?description=vatiala", []Line{
			{lineUrl("1"), "1", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}},
		{"/lines?name=1&exclude-fields=name,description", []Line{
			{lineUrl("1"), "", ""},
		}},
		{"/lines?name=1&exclude-fields=description", []Line{
			{lineUrl("1"), "1", ""},
		}},
		{"/lines?exclude-fields=name", []Line{
			{lineUrl("-1"), "", "Foobar"},
			{lineUrl("1"), "", "Vatiala - Pirkkala"},
			{lineUrl("1A"), "", "Vatiala - Pirkkala (lentoasema)"},
			{lineUrl("3A"), "", "Etelä-Hervanta - Lentävänniemi"},
		}},
		{"/lines/1A", []Line{
			{lineUrl("1A"), "1A", "Vatiala - Pirkkala (lentoasema)"},
		}},
		{"/lines/foobar", []Line{}},
		{"/lines?name=noSuchThing", []Line{}},
		{"/lines?description=noSuchThing", []Line{}},
		{"/lines?description=noSuchThing&name=noSuchThing", []Line{}},
		{"/lines?description=vatiala&name=noSuchThing", []Line{}},
		{"/lines?description=noSuchThing&name=1", []Line{}},
	}

	for _, tc := range testCases {
		router, recorder, ctx := initializeTest(t)
		InjectLineRoutes(router, ctx)

		gotResponse := triggerGetLinesRequest(t, router, recorder, tc.url)
		expectedResponse := lineSuccessResponse{
			Status: "success",
			Data: apiSuccessData{
				Headers: apiHeaders{
					Paging: apiHeadersPaging{
						StartIndex: 0,
						PageSize:   uint16(len(tc.items)),
						MoreData:   false,
					},
				},
			},
			Body: tc.items,
		}

		var diffs []FieldDiff
		initialTag := fmt.Sprintf("%v:Response", tc.url)
		var err = compareVariables(expectedResponse, gotResponse, initialTag, &diffs, false)
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

type lineSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Line         `json:"body"`
}

func triggerGetLinesRequest(t *testing.T, router *mux.Router, recorder *httptest.ResponseRecorder, target string) lineSuccessResponse {
	serveHttp(t, router, recorder, target)

	var response lineSuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}
