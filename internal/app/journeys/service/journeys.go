package service

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"strings"
	"time"
)

type JourneysService struct {
	Repository *repository.JourneysRepository
}

func (s JourneysService) Search(params map[string]string, excludeInactive bool) []*model.Journey {
	result := make([]*model.Journey, 0)

	for _, journey := range s.Repository.Journeys.All {
		if journeyMatchesConditions(journey, params, excludeInactive) {
			result = append(result, journey)
		}
	}

	return result
}

func (s JourneysService) GetOneById(id string) (*model.Journey, error) {
	var journey *model.Journey

	if j, ok := s.Repository.Journeys.ById[id]; ok {
		journey = j
	}

	if j, ok := s.Repository.Journeys.ByActivityId[id]; ok {
		journey = j
	}

	if journey != nil {
		return journey, nil
	}

	return nil, model.ErrNoSuchElement
}

func journeyMatchesConditions(journey *model.Journey, conditions map[string]string, excludeInactive bool) bool {
	now := time.Now()
	curDay := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

	if journey == nil || (excludeInactive && !(journey.ValidFrom <= curDay && journey.ValidTo >= curDay)) {
		return false
	}

	if conditions == nil {
		return true
	}

	for k, v := range conditions {
		switch k {
		case "lineId":
			if journey.Line == nil || journey.Line.Name != v {
				return false
			}
		case "routeId":
			if journey.Route == nil || journey.Route.Id != v {
				return false
			}
		case "journeyPatternId":
			if journey.JourneyPattern == nil || journey.JourneyPattern.Id != v {
				return false
			}
		case "dayTypes":
			matched := false
			vDayTypes := strings.Split(v, ",")
			for _, dt := range journey.DayTypes {
				for _, vdt := range vDayTypes {
					if dt == vdt {
						matched = true
						break
					}
				}
			}
			if !matched {
				return false
			}
		case "departureTime":
			// Check also with ":00" postfix to be backwards compatible with the old API
			if journey.DepartureTime != v && fmt.Sprintf("%v:00", v) != journey.DepartureTime {
				return false
			}
		case "arrivalTime":
			if journey.ArrivalTime != v && fmt.Sprintf("%v:00", v) != journey.ArrivalTime {
				return false
			}
		case "firstStopPointId":
			if len(journey.Calls) == 0 {
				return false
			}
			first := journey.Calls[0]
			if first.StopPoint == nil || first.StopPoint.ShortName != v {
				return false
			}
		case "lastStopPointId":
			if len(journey.Calls) == 0 {
				return false
			}
			last := journey.Calls[len(journey.Calls)-1]
			if last.StopPoint == nil || last.StopPoint.ShortName != v {
				return false
			}
		case "stopPointId":
			matched := false
			for _, c := range journey.Calls {
				if c.StopPoint != nil && c.StopPoint.ShortName == v {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		case "gtfsTripId":
			if journey.GtfsInfo == nil || journey.GtfsInfo.TripId != v {
				return false
			}
		}
	}

	return true
}
