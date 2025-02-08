package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"log"
	"sort"
	"strconv"
	"strings"
)

func newRoutesRepository(shapes []*ggtfs.Shape) *JourneysRoutesRepository {
	var all = make([]*model.Route, 0)
	var byId = make(map[string]*model.Route)

	var shapeIdToCoords = make(map[string][][]float64)
	for i, shape := range shapes {
		if shape == nil {
			fmt.Println(fmt.Sprintf("Nil shape detected, number %v in the shapes array, newRoutesRepository function", i))
			continue
		}

		if shape.Id == nil {
			fmt.Println(fmt.Sprintf("Shape.Id is missing, GTFS line: %v", shape.LineNumber))
			continue
		}

		shapeId := strings.TrimSpace(*shape.Id)

		var lat, lon float64
		if shape.PtLat != nil {
			lf, err := strconv.ParseFloat(*shape.PtLat, 64)
			if err != nil {
				log.Println(fmt.Sprintf("shape (on gtfs line %v): cannot parse lat float value", shape.LineNumber))
			}
			lat = lf
		} else {
			log.Println(fmt.Sprintf("shape (on gtfs line %v): lat is missing", shape.LineNumber))
		}

		if shape.PtLon != nil {
			lf, err := strconv.ParseFloat(*shape.PtLon, 64)
			if err != nil {
				log.Println(fmt.Sprintf("shape (on gtfs line %v): cannot parse lon float value", shape.LineNumber))
			}
			lon = lf
		} else {
			log.Println(fmt.Sprintf("shape (on gtfs line %v): lon is missing", shape.LineNumber))
		}

		if _, ok := shapeIdToCoords[shapeId]; !ok {
			shapeIdToCoords[shapeId] = make([][]float64, 0)
		}

		shapeIdToCoords[shapeId] = append(shapeIdToCoords[shapeId], []float64{lat, lon})
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

	return &JourneysRoutesRepository{
		All:  all,
		ById: byId,
	}
}

func createCoordinateProjection(coords [][]float64) (string, error) {
	if coords == nil || len(coords) == 0 {
		return "", nil
	}

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

type JourneysRoutesRepository struct {
	All  []*model.Route
	ById map[string]*model.Route
}
