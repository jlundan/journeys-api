package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
)

type RoutesService struct {
	DataStore *repository.JourneysDataStore
}

func (s RoutesService) Search(params map[string]string) []*model.Route {
	result := make([]*model.Route, 0)

	for _, route := range s.DataStore.Routes.All {
		if routeMatchesConditions(route, params) {
			result = append(result, route)
		}
	}

	return result
}

func (s RoutesService) GetOneById(id string) (*model.Route, error) {
	if r, ok := s.DataStore.Routes.ById[id]; ok {
		return r, nil
	}
	return nil, model.ErrNoSuchElement
}

func routeMatchesConditions(route *model.Route, conditions map[string]string) bool {
	if route == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(route.Name, v) {
				return false
			}
		case "lineId":
			if route.Line == nil || route.Line.Name != v {
				return false
			}
		}
	}

	return true
}
