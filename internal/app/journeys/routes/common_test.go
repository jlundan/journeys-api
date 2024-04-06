package routes

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/context/tre"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockFilterObject, originalFilterObject func(r interface{}, properties string) (interface{}, error)
var fakeMarshal, originalMarshal func(_ interface{}) ([]byte, error)

func init() {
	_ = os.Setenv("JOURNEYS_BASE_URL", "http://localhost:5678")
	_ = os.Setenv("JOURNEYS_GTFS_PATH", "testdata/tre1/gtfs")

	mockFilterObject = func(r interface{}, properties string) (interface{}, error) {
		return nil, errors.New("foo")
	}

	fakeMarshal = func(_ interface{}) ([]byte, error) {
		return []byte{}, errors.New("marshalling failed")
	}
}

func TestCreateSuccessResponseFail(t *testing.T) {
	originalFilterObject = filterObject
	filterObject = mockFilterObject
	_, err := createSuccessResponse(make([]interface{}, 1), 0, 0, false, "")
	if err == nil {
		t.Error("Expected to see an error")
	}

	filterObject = originalFilterObject
}

func TestSendResponseFail(t *testing.T) {
	w := &responseWriterMock{}
	sendResponse(nil, errors.New("foobar"), nil, w)
	if string(w.BodyData) != `{"status":"fail","data":{"message":"foobar"}}` {
		t.Error("expected an error message")
	}

	w.Reset()

	originalFilterObject = filterObject
	filterObject = mockFilterObject
	sendResponse(make([]interface{}, 1), nil, nil, w)
	if string(w.BodyData) != `{"status":"fail","data":{"message":"foo"}}` {
		t.Error("expected an error message")
	}
	filterObject = originalFilterObject

	w.Reset()

	originalMarshal = jsonMarshal
	jsonMarshal = fakeMarshal
	sendResponse(nil, nil, nil, w)
	if w.StatusCode != 500 {
		t.Error("expected status code: 500")
	}
	jsonMarshal = originalMarshal

	w.Reset()
	w.ErrorOnWrite = true

	sendResponse(nil, nil, nil, w)
	if w.StatusCode != 500 {
		t.Error("expected status code: 500")
	}
	w.Reset()

}

func validateCommonResponseFields(t *testing.T, status string, data apiSuccessData, dataSize uint16) bool {
	if status != "success" {
		t.Errorf("expected response status to be success, got %v", status)
		return false
	}

	if data.Headers.Paging.StartIndex != 0 {
		t.Errorf("expected response startIndex to be 0, got %v", data.Headers.Paging.StartIndex)
		return false
	}

	if data.Headers.Paging.PageSize != dataSize {
		t.Errorf("expected response data size to be %v, got %v", dataSize, data.Headers.Paging.PageSize)
		return false
	}

	if data.Headers.Paging.MoreData != false {
		t.Errorf("expected response moredata to be false, got %v", data.Headers.Paging.MoreData)
		return false
	}

	return true
}

func initializeTest(t *testing.T) (*mux.Router, *httptest.ResponseRecorder, model.Context) {
	ctx, err, _ := tre.NewContext()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	r := mux.NewRouter()

	return r, w, ctx
}

func serveHttp(t *testing.T, r *mux.Router, w *httptest.ResponseRecorder, target string) {
	r.ServeHTTP(w, httptest.NewRequest("GET", target, nil))

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
}

func getStopPointMap() map[string]StopPoint {
	result := make(map[string]StopPoint)

	stopPoints := []struct {
		id           string
		name         string
		location     string
		tariffZone   string
		municipality Municipality
	}{
		{"4600", "Vatiala", "61.47561,23.97756", "B", getMunicipalityMap()["211"]},
		{"8171", "Vällintie", "61.48067,23.97002", "B", getMunicipalityMap()["211"]},
		{"8149", "Sudenkorennontie", "61.47979,23.96166", "C", getMunicipalityMap()["211"]},
		{"7017", "Suupantori", "61.46546,23.64219", "B", getMunicipalityMap()["604"]},
		{"7015", "Pirkkala", "61.4659,23.64734", "B", getMunicipalityMap()["604"]},
		{"3615", "Näyttelijänkatu", "61.4445,23.87235", "B", getMunicipalityMap()["837"]},
		{"3607", "Lavastajanpolku", "61.44173,23.86961", "B", getMunicipalityMap()["837"]},
	}

	for _, tc := range stopPoints {
		result[tc.id] = StopPoint{
			Url:          stopPointUrl(tc.id),
			ShortName:    tc.id,
			Name:         tc.name,
			Location:     tc.location,
			TariffZone:   tc.tariffZone,
			Municipality: tc.municipality,
		}
	}

	return result
}

func getMunicipalityMap() map[string]Municipality {
	municipalities := make(map[string]Municipality)
	municipalities["211"] = Municipality{
		Url:       municipalityUrl("211"),
		ShortName: "211",
		Name:      "Kangasala",
	}

	municipalities["604"] = Municipality{
		Url:       municipalityUrl("604"),
		ShortName: "604",
		Name:      "Pirkkala",
	}

	municipalities["837"] = Municipality{
		Url:       municipalityUrl("837"),
		ShortName: "837",
		Name:      "Tampere",
	}
	return municipalities
}

func getJourneyPatternMap() map[string]JourneyPattern {
	result := make(map[string]JourneyPattern)

	stopPoints := []struct {
		id              string
		line            string
		route           string
		originStop      string
		destinationStop string
		name            string
		direction       string
		stopPoints      []StopPoint
	}{
		{"047b0afc973ee2fd4fe92b128c3a932a", "1", "1504270174600", "7017",
			"7015", "Suupantori - Pirkkala", "1",
			[]StopPoint{getStopPointMap()["7017"], getStopPointMap()["7015"]}},

		{"65f51d2f85284af2fad1305c0ce71033", "3A", "1517136151028", "3615",
			"3607", "Näyttelijänkatu - Lavastajanpolku", "0",
			[]StopPoint{getStopPointMap()["3615"], getStopPointMap()["3607"]}},

		{"9bc7403ad27267edbfbd63c3e92e5afa", "1A", "1501146007035", "4600",
			"8149", "Vatiala - Sudenkorennontie", "0",
			[]StopPoint{getStopPointMap()["4600"], getStopPointMap()["8171"], getStopPointMap()["8149"]}},

		{"c01c71b0c9f456ba21f498a1dca54b3b", "-1", "111111111", "3615",
			"7017", "Näyttelijänkatu - Suupantori", "0",
			[]StopPoint{getStopPointMap()["3615"], getStopPointMap()["7017"]}},
	}

	for _, tc := range stopPoints {
		result[tc.id] = JourneyPattern{
			Url:             journeyPatternUrl(tc.id),
			LineUrl:         lineUrl(tc.line),
			RouteUrl:        routeUrl(tc.route),
			OriginStop:      stopPointUrl(tc.originStop),
			DestinationStop: stopPointUrl(tc.destinationStop),
			Direction:       tc.direction,
			Name:            tc.name,
			StopPoints:      tc.stopPoints,
		}
	}

	return result
}

func getJourneyMap() map[string]Journey {
	result := make(map[string]Journey)

	journeys := []struct {
		id                   string
		line                 string
		activityId           string
		route                string
		journeyPattern       string
		departureTime        string
		arrivalTime          string
		headSign             string
		directionId          string
		wheelchairAccessible bool
		gtfs                 JourneyGtfsInfo
		dayTypes             []string
		dayTypeExceptions    []DayTypeException
		calls                []JourneyCall
	}{
		{
			"111111111",
			"-1",
			"-1_0720_7017_3615",
			"111111111",
			"c01c71b0c9f456ba21f498a1dca54b3b",
			"07:20:00",
			"07:21:00",
			"Foobar",
			"0",
			false,
			JourneyGtfsInfo{TripId: "111111111"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{},
			[]JourneyCall{
				{"07:20:00", "07:20:00", getStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getStopPointMap()["7017"]},
			},
		},
		{
			"7020205685",
			"1",
			"1_1443_7015_7017",
			"1504270174600",
			"047b0afc973ee2fd4fe92b128c3a932a",
			"14:43:00",
			"14:44:45",
			"Vatiala",
			"1",
			false,
			JourneyGtfsInfo{TripId: "7020205685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"14:43:00", "14:43:00", getStopPointMap()["7017"]},
				{"14:44:45", "14:44:45", getStopPointMap()["7015"]},
			},
		},
		{
			"7020295685",
			"1A",
			"1A_0630_8149_4600",
			"1501146007035",
			"9bc7403ad27267edbfbd63c3e92e5afa",
			"06:30:00",
			"06:32:30",
			"Lentoasema",
			"0",
			false,
			JourneyGtfsInfo{TripId: "7020295685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"06:30:00", "06:30:00", getStopPointMap()["4600"]},
				{"06:31:30", "06:31:30", getStopPointMap()["8171"]},
				{"06:32:30", "06:32:30", getStopPointMap()["8149"]},
			},
		},
		{
			"7024545685",
			"3A",
			"3A_0720_3607_3615",
			"1517136151028",
			"65f51d2f85284af2fad1305c0ce71033",
			"07:20:00",
			"07:21:00",
			"Lentävänniemi",
			"0",
			false,
			JourneyGtfsInfo{TripId: "7024545685"},
			[]string{"monday", "tuesday", "wednesday", "thursday", "friday"},
			[]DayTypeException{{"2021-04-05", "2021-04-05", "yes"}, {"2021-05-13", "2021-05-13", "no"}},
			[]JourneyCall{
				{"07:20:00", "07:20:00", getStopPointMap()["3615"]},
				{"07:21:00", "07:21:00", getStopPointMap()["3607"]},
			},
		},
	}

	for _, tc := range journeys {
		result[tc.id] = Journey{
			Url:                  journeyUrl(tc.id),
			ActivityUrl:          journeyActivityUrl(tc.activityId),
			LineUrl:              lineUrl(tc.line),
			RouteUrl:             routeUrl(tc.route),
			JourneyPatternUrl:    journeyPatternUrl(tc.journeyPattern),
			DepartureTime:        tc.departureTime,
			ArrivalTime:          tc.arrivalTime,
			HeadSign:             tc.headSign,
			Direction:            tc.directionId,
			WheelchairAccessible: tc.wheelchairAccessible,
			GtfsInfo:             tc.gtfs,
			DayTypes:             tc.dayTypes,
			DayTypeExceptions:    tc.dayTypeExceptions,
			Calls:                tc.calls,
		}
	}

	return result
}

func journeysMatch(a Journey, b Journey) bool {
	if a.Url != b.Url || a.LineUrl != b.LineUrl || a.RouteUrl != b.RouteUrl || a.JourneyPatternUrl != b.JourneyPatternUrl ||
		a.DepartureTime != b.DepartureTime || a.ArrivalTime != b.ArrivalTime || a.HeadSign != b.HeadSign ||
		a.WheelchairAccessible != b.WheelchairAccessible || a.GtfsInfo != b.GtfsInfo || len(a.DayTypes) != len(b.DayTypes) ||
		a.ActivityUrl != b.ActivityUrl || len(a.DayTypeExceptions) != len(b.DayTypeExceptions) || len(a.Calls) != len(b.Calls) {
		return false
	}

	for i := range a.DayTypes {
		if a.DayTypes[i] != b.DayTypes[i] {
			return false
		}
	}

	for i := range a.DayTypeExceptions {
		if a.DayTypeExceptions[i] != b.DayTypeExceptions[i] {
			return false
		}
	}

	for i := range a.Calls {
		if a.Calls[i] != b.Calls[i] {
			return false
		}
	}

	return true
}

func lineUrl(name string) string {
	return fmt.Sprintf("%v/lines/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func journeyPatternUrl(name string) string {
	return fmt.Sprintf("%v/journey-patterns/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func journeyUrl(name string) string {
	return fmt.Sprintf("%v/journeys/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func municipalityUrl(name string) string {
	return fmt.Sprintf("%v/municipalities/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func routeUrl(name string) string {
	return fmt.Sprintf("%v/routes/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func stopPointUrl(name string) string {
	return fmt.Sprintf("%v/stop-points/%v", os.Getenv("JOURNEYS_BASE_URL"), name)
}

func journeyActivityUrl(name string) string {
	return fmt.Sprintf("%v/vehicle-activity/%v", os.Getenv("JOURNEYS_VA_BASE_URL"), name)
}

type responseWriterMock struct {
	BodyData     []byte
	StatusCode   int
	ErrorOnWrite bool
}

func (rw *responseWriterMock) Header() http.Header {
	return http.Header{}
}

func (rw *responseWriterMock) Write(bytes []byte) (int, error) {
	if rw.ErrorOnWrite {
		return 0, errors.New("write error")
	}

	rw.BodyData = bytes
	return len(bytes), nil
}

func (rw *responseWriterMock) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func (rw *responseWriterMock) Reset() {
	rw.StatusCode = 0
	rw.BodyData = nil
	rw.ErrorOnWrite = false
}
