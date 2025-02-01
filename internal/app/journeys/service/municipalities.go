package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
)

type MunicipalitiesService struct {
	DataStore *repository.JourneysRepository
}

func (s MunicipalitiesService) Search(params map[string]string) []*model.Municipality {
	result := make([]*model.Municipality, 0)

	for _, municipality := range s.DataStore.Municipalities.All {
		if municipalityMatchesConditions(municipality, params) {
			result = append(result, municipality)
		}
	}

	return result
}

func (s MunicipalitiesService) GetOneById(id string) (*model.Municipality, error) {
	if m, ok := s.DataStore.Municipalities.ById[id]; ok {
		return m, nil
	}
	return nil, model.ErrNoSuchElement
}

func municipalityMatchesConditions(municipality *model.Municipality, conditions map[string]string) bool {
	if municipality == nil {
		return false
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if !utils.StrContains(municipality.Name, v) {
				return false
			}
		case "shortName":
			if !utils.StrContains(municipality.PublicCode, v) {
				return false
			}
		}
	}

	return true
}
