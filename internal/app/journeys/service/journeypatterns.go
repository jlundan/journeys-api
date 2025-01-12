package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
)

func (ds DataService) GetOneJourneyPatternById(id string) (*model.JourneyPattern, error) {
	if jp, ok := ds.DataStore.JourneyPatterns.ById[id]; ok {
		return jp, nil
	}
	return nil, model.ErrNoSuchElement
}

func (ds DataService) SearchJourneyPatterns(params map[string]string) []*model.JourneyPattern {
	result := make([]*model.JourneyPattern, 0)

	for _, jp := range ds.DataStore.JourneyPatterns.All {
		if journeyPatternMatchesConditions(jp, params) {
			result = append(result, jp)
		}
	}

	return result
}

func journeyPatternMatchesConditions(journeyPattern *model.JourneyPattern, conditions map[string]string) bool {
	if journeyPattern == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(journeyPattern.Name, v) {
				return false
			}
		case "lineId":
			if journeyPattern.Route == nil || journeyPattern.Route.Line == nil || journeyPattern.Route.Line.Name != v {
				return false
			}
		case "firstStopPointId":
			if len(journeyPattern.StopPoints) == 0 || journeyPattern.StopPoints[0].ShortName != v {
				return false
			}
		case "lastStopPointId":
			spLength := len(journeyPattern.StopPoints)
			if spLength == 0 || journeyPattern.StopPoints[spLength-1].ShortName != v {
				return false
			}
		case "stopPointId":
			found := false
			for _, sp := range journeyPattern.StopPoints {
				if sp.ShortName == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}
