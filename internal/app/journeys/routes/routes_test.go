package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestRoutesMatchesConditions(t *testing.T) {
	testCases := []struct {
		item       *model.Route
		conditions map[string]string
		expected   bool
	}{
		{nil, nil, false},
		{&model.Route{Line: nil}, map[string]string{"lineId": "1"}, false},
	}

	for _, tc := range testCases {
		matches := routeMatchesConditions(tc.item, tc.conditions)
		if matches != tc.expected {
			t.Error(fmt.Sprintf("expected %v, got %v", matches, tc.expected))
		}
	}
}

func TestGetRoutes(t *testing.T) {
	rm := getRouteMap()

	// Route id is GTFS shape ID
	testCases := []struct {
		target   string
		expected []Route
	}{
		{"/routes",
			[]Route{rm["111111111"], rm["1501146007035"], rm["1504270174600"], rm["1517136151028"]},
		},
		{"/routes/1501146007035",
			[]Route{rm["1501146007035"]},
		},
		{"/routes?lineId=1A",
			[]Route{rm["1501146007035"]},
		},
		{"/routes/foobar",
			[]Route{},
		},
	}

	for _, tc := range testCases {
		r, w, ctx := initializeTest(t)
		InjectRouteRoutes(r, ctx)

		routesResponse := sendRoutesRequest(t, r, w, tc.target)

		dataSize := len(tc.expected)
		if success := validateCommonResponseFields(t, routesResponse.Status, routesResponse.Data, uint16(dataSize)); !success {
			break
		}
		if len(routesResponse.Body) != dataSize {
			t.Errorf("expected %v, got %v", dataSize, len(routesResponse.Body))
			break
		}
		for i, route := range routesResponse.Body {
			if tc.expected[i].Projection != route.Projection || !reflect.DeepEqual(route.Links, tc.expected[i].Links) {
				t.Errorf("expected %v, got %v", tc.expected[i], route)
				break
			}
		}
	}
}

type routeSuccessResponse struct {
	Status string         `json:"status"`
	Data   apiSuccessData `json:"data"`
	Body   []Route        `json:"body"`
}

func sendRoutesRequest(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) routeSuccessResponse {
	serveHttp(t, r, w, target)

	var response routeSuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	return response
}

func getRouteMap() map[string]Route {
	result := make(map[string]Route)

	routes := []struct {
		id         string
		projection string
		links      map[string]string
	}{
		{
			"111111111",
			"6144444,2387262:17,-14",
			map[string]string{
				"self":             routeUrl("111111111"),
				"line":             lineUrl("-1"),
				"journey-patterns": fmt.Sprintf("%v/journey-patterns?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "111111111"),
				"journeys":         fmt.Sprintf("%v/journeys?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "111111111"),
			},
		},
		{
			"1501146007035",
			"6147590,2397764:-104,32",
			map[string]string{
				"self":             routeUrl("1501146007035"),
				"line":             lineUrl("1A"),
				"journey-patterns": fmt.Sprintf("%v/journey-patterns?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1501146007035"),
				"journeys":         fmt.Sprintf("%v/journeys?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1501146007035"),
			},
		},
		{
			"1504270174600",
			"6146557,2364305:-19,-148:-17,-191",
			map[string]string{
				"self":             routeUrl("1504270174600"),
				"line":             lineUrl("1"),
				"journey-patterns": fmt.Sprintf("%v/journey-patterns?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1504270174600"),
				"journeys":         fmt.Sprintf("%v/journeys?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1504270174600"),
			},
		},
		{
			"1517136151028",
			"6144444,2387262:17,-14",
			map[string]string{
				"self":             routeUrl("1517136151028"),
				"line":             lineUrl("3A"),
				"journey-patterns": fmt.Sprintf("%v/journey-patterns?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1517136151028"),
				"journeys":         fmt.Sprintf("%v/journeys?routeId=%v", os.Getenv("JOURNEYS_BASE_URL"), "1517136151028"),
			},
		},
	}

	for _, route := range routes {
		result[route.id] = Route{
			Projection: route.projection,
			Links:      route.links,
		}
	}

	return result
}
