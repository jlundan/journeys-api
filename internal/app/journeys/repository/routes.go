package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"sort"
	"strconv"
)

type JourneysRouteDataStore struct {
	All  []*model.Route
	ById map[string]*model.Route
}

func newRouteDataStore(shapes []*ggtfs.Shape) *JourneysRouteDataStore {
	var all = make([]*model.Route, 0)
	var byId = make(map[string]*model.Route)

	var shapeIdToCoords = make(map[string][][]float64)
	for _, shape := range shapes {
		// TODO: nil check
		if _, ok := shapeIdToCoords[*shape.Id]; !ok {
			shapeIdToCoords[*shape.Id] = make([][]float64, 0)
		}
		lat, err := strconv.ParseFloat(*shape.PtLat, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error parsing shape.PtLat: %v, line: %v", shape.PtLat, shape.LineNumber))
		}
		lon, err := strconv.ParseFloat(*shape.PtLon, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error parsing shape.PtLon: %v, line: %v", shape.PtLon, shape.LineNumber))
		}

		shapeIdToCoords[*shape.Id] = append(shapeIdToCoords[*shape.Id], []float64{lat, lon})
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

	return &JourneysRouteDataStore{
		All:  all,
		ById: byId,
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
