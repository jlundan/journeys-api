package tre

import (
	"errors"
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"sort"
)

type Routes struct {
	all  []*model.Route
	byId map[string]*model.Route
}

func (routes Routes) GetOne(name string) (*model.Route, error) {
	if _, ok := routes.byId[name]; !ok {
		return &model.Route{}, errors.New("no such element")
	}
	return routes.byId[name], nil
}
func (routes Routes) GetAll() []*model.Route {
	return routes.all
}

func buildRoutes(g GTFSContext) Routes {
	var all = make([]*model.Route, 0)
	var byId = make(map[string]*model.Route)

	var shapeIdToCoords = make(map[string][][]float64)
	for _, shape := range g.Shapes {
		if _, ok := shapeIdToCoords[shape.Id]; !ok {
			shapeIdToCoords[shape.Id] = make([][]float64, 0)
		}
		shapeIdToCoords[shape.Id] = append(shapeIdToCoords[shape.Id], []float64{shape.PtLat, shape.PtLon})
	}

	for shapeId, coords := range shapeIdToCoords {
		projection, err := createCoordinateProjection(coords)
		if err != nil {
			continue
		}
		route := &model.Route{
			Id:            shapeId,
			GeoProjection: projection,
			// Line, Name , JourneyPatterns and Journeys are populated while building journeys
		}

		all = append(all, route)
		byId[shapeId] = route
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Id < all[y].Id
	})

	return Routes{
		all:  all,
		byId: byId,
	}
}

func createCoordinateProjection(coords [][]float64) (string, error) {
	projection := ""
	var lastLat int64
	var lastLon int64

	for k, v := range coords {
		lat := int64(v[0] * 100000)
		lon := int64(v[1] * 100000)
		if k == 0 {
			projection += fmt.Sprintf("%v,%v", lat, lon)
		} else {
			projection += fmt.Sprintf(":%v,%v", lastLat-lat, lastLon-lon)
		}

		lastLat = lat
		lastLon = lon
	}
	return projection, nil
}
