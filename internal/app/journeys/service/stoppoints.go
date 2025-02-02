package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
	"strconv"
	"strings"
)

type StopPointsService struct {
	Repository *repository.JourneysRepository
}

func (s StopPointsService) Search(params map[string]string) []*model.StopPoint {
	result := make([]*model.StopPoint, 0)

	for _, sp := range s.Repository.StopPoints.All {
		if stopPointMatchesConditions(sp, params) {
			result = append(result, sp)
		}
	}

	return result
}

func (s StopPointsService) GetOneById(id string) (*model.StopPoint, error) {
	if sp, ok := s.Repository.StopPoints.ById[id]; ok {
		return sp, nil
	}
	return nil, model.ErrNoSuchElement
}

func stopPointMatchesConditions(stopPoint *model.StopPoint, conditions map[string]string) bool {
	if stopPoint == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(stopPoint.Name, v) {
				return false
			}
		case "shortName":
			if !utils.StrContains(stopPoint.ShortName, v) {
				return false
			}
		case "tariffZone":
			if !utils.StrContains(stopPoint.TariffZone, v) {
				return false
			}
		case "municipalityName":
			if stopPoint.Municipality == nil || !utils.StrContains(stopPoint.Municipality.Name, v) {
				return false
			}
		case "municipalityShortName":
			if stopPoint.Municipality == nil || !utils.StrContains(stopPoint.Municipality.PublicCode, v) {
				return false
			}
		case "location":
			return stopPointLocationMatches(stopPoint, v)
		}
	}

	return true
}

func stopPointLocationMatches(stopPoint *model.StopPoint, locationExpression string) bool {
	locationParts := strings.Split(locationExpression, ":")
	if len(locationParts) == 2 {
		upperLeft := strings.Split(locationParts[0], ",")
		if len(upperLeft) != 2 {
			return false
		}

		upperLeftLat, err := strconv.ParseFloat(upperLeft[0], 64)
		if err != nil {
			return false
		}

		upperLeftLon, err := strconv.ParseFloat(upperLeft[1], 64)
		if err != nil {
			return false
		}

		lowerRight := strings.Split(locationParts[1], ",")
		if len(lowerRight) != 2 {
			return false
		}

		lowerRightLat, err := strconv.ParseFloat(lowerRight[0], 64)
		if err != nil {
			return false
		}

		lowerRightLon, err := strconv.ParseFloat(lowerRight[1], 64)
		if err != nil {
			return false
		}

		return upperLeftLat <= stopPoint.Latitude && stopPoint.Latitude <= lowerRightLat &&
			upperLeftLon <= stopPoint.Longitude && stopPoint.Longitude <= lowerRightLon

	}

	return fmt.Sprintf("%v,%v", stopPoint.Latitude, stopPoint.Longitude) == locationExpression
}
