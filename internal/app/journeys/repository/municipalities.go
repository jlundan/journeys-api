package repository

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"sort"
)

func newMunicipalitiesRepository(m municipalityData) *JourneysMunicipalitiesRepository {
	var all = make([]*model.Municipality, 0)
	var byId = make(map[string]*model.Municipality)

	for _, v := range m.municipalityRows {
		municipality := model.Municipality{
			PublicCode: v[m.municipalityHeaders["id"]],
			Name:       v[m.municipalityHeaders["name"]],
		}
		all = append(all, &municipality)
		byId[v[m.municipalityHeaders["id"]]] = &municipality
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].PublicCode < all[y].PublicCode
	})

	return &JourneysMunicipalitiesRepository{
		All:  all,
		ById: byId,
	}
}

type JourneysMunicipalitiesRepository struct {
	All  []*model.Municipality
	ById map[string]*model.Municipality
}
