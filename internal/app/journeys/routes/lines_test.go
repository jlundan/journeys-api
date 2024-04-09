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
		target string
		items  []Line
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
		r, w, ctx := initializeTest(t)
		InjectLineRoutes(r, ctx)

		linesResponse := triggerGetLinesRequest(t, r, w, tc.target)

		dataSize := len(tc.items)
		if success := validateCommonResponseFields(t, linesResponse.Status, linesResponse.Data, uint16(dataSize)); !success {
			break
		}
		if len(linesResponse.Body) != dataSize {
			t.Errorf("expected %v, got %v", dataSize, len(linesResponse.Body))
			break
		}
		for i, line := range linesResponse.Body {
			if tc.items[i] != line {
				t.Errorf("expected %v, got %v", tc.items[i], line)
				break
			}
		}
	}
}

type lineSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Line         `json:"body"`
}

func triggerGetLinesRequest(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) lineSuccessResponse {
	serveHttp(t, r, w, target)

	var response lineSuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}
