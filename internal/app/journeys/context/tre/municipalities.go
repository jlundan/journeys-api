package tre

import (
	"errors"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"sort"
)

type Municipalities struct {
	all  []*model.Municipality
	byId map[string]*model.Municipality
}

func (municipalities Municipalities) GetOne(name string) (*model.Municipality, error) {
	if _, ok := municipalities.byId[name]; !ok {
		return nil, errors.New("no such element")
	}
	return municipalities.byId[name], nil
}
func (municipalities Municipalities) GetAll() []*model.Municipality {
	return municipalities.all
}

func buildMunicipalities(m municipalityData) Municipalities {
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

	return Municipalities{
		all:  all,
		byId: byId,
	}
}
