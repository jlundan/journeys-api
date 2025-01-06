package routes

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const stopPointPrefix = "/stop-points"

func InjectStopPointRoutes(r *mux.Router, context model.Context) {
	sr := r.PathPrefix(stopPointPrefix).Subrouter()

	sr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllStopPoints(w, r, context)
	}).Methods("GET")

	sr.HandleFunc(`/{name}`, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		handleGetOneStopPoint(w, r, context, params["name"])
	}).Methods("GET")
}

func handleGetAllStopPoints(w http.ResponseWriter, r *http.Request, context model.Context) {
	responseItems := make([]interface{}, 0)

	for _, stopPoint := range context.StopPoints().GetAll() {
		matches, err := stopPointMatchesConditions(stopPoint, getDefaultConditions(r))
		if err != nil {
			sendError(err, w)
			return
		}
		if matches {
			responseItems = append(responseItems, convertStopPoint(stopPoint))
		}
	}

	sendResponse(responseItems, nil, r, w)
}

func handleGetOneStopPoint(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
	stopPoint, err := context.StopPoints().GetOne(name)

	if err != nil {
		sendResponse(nil, err, r, w)
	} else {
		sendResponse([]interface{}{convertStopPoint(stopPoint)}, nil, r, w)
	}
}

func convertStopPoint(stopPoint *model.StopPoint) StopPoint {
	return StopPoint{
		Url:          fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), stopPointPrefix, stopPoint.ShortName),
		ShortName:    stopPoint.ShortName,
		Name:         stopPoint.Name,
		Location:     fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude),
		TariffZone:   stopPoint.TariffZone,
		Municipality: convertMunicipality(stopPoint.Municipality),
	}
}

func stopPointMatchesConditions(stopPoint *model.StopPoint, conditions map[string]string) (bool, error) {
	if stopPoint == nil {
		return false, nil
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(stopPoint.Name, v) {
				return false, nil
			}
		case "shortName":
			if !utils.StrContains(stopPoint.ShortName, v) {
				return false, nil
			}
		case "tariffZone":
			if !utils.StrContains(stopPoint.TariffZone, v) {
				return false, nil
			}
		case "municipalityName":
			if stopPoint.Municipality == nil || !utils.StrContains(stopPoint.Municipality.Name, v) {
				return false, nil
			}
		case "municipalityShortName":
			if stopPoint.Municipality == nil || !utils.StrContains(stopPoint.Municipality.PublicCode, v) {
				return false, nil
			}
		case "location":
			locationMatches, err := stopPointLocationMatches(stopPoint, v)
			if err != nil {
				return false, err
			}
			if !locationMatches {
				return false, nil
			}
		}
	}

	return true, nil
}

func stopPointLocationMatches(stopPoint *model.StopPoint, locationExpression string) (bool, error) {
	locationParts := strings.Split(locationExpression, ":")
	if len(locationParts) == 2 {
		upperLeft := strings.Split(locationParts[0], ",")
		if len(upperLeft) != 2 {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		upperLeftLat, err := strconv.ParseFloat(upperLeft[0], 64)
		if err != nil {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		upperLeftLon, err := strconv.ParseFloat(upperLeft[1], 64)
		if err != nil {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		lowerRight := strings.Split(locationParts[1], ",")
		if len(lowerRight) != 2 {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		lowerRightLat, err := strconv.ParseFloat(lowerRight[0], 64)
		if err != nil {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		lowerRightLon, err := strconv.ParseFloat(lowerRight[1], 64)
		if err != nil {
			return false, errors.New("illegal coordinate format: must be latitude,longitude")
		}

		return upperLeftLat <= stopPoint.Latitude && stopPoint.Latitude <= lowerRightLat &&
			upperLeftLon <= stopPoint.Longitude && stopPoint.Longitude <= lowerRightLon, nil

	}

	return fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude) == locationExpression, nil
}

type StopPoint struct {
	Url          string       `json:"url"`
	ShortName    string       `json:"shortName"`
	Name         string       `json:"name"`
	Location     string       `json:"location"`
	TariffZone   string       `json:"tariffZone"`
	Municipality Municipality `json:"municipality"`
}
