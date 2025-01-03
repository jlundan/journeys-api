package tre

import (
	"errors"
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
	"math"
	"sort"
	"strconv"
)

type StopPoints struct {
	all  []*model.StopPoint
	byId map[string]*model.StopPoint
}

func (stopPoints StopPoints) GetOne(name string) (*model.StopPoint, error) {
	if _, ok := stopPoints.byId[name]; !ok {
		return &model.StopPoint{}, errors.New("no such element")
	}
	return stopPoints.byId[name], nil
}
func (stopPoints StopPoints) GetAll() []*model.StopPoint {
	return stopPoints.all
}

func buildStopPoints(g GTFSContext, municipalities Municipalities) StopPoints {
	var warnings []error

	var all = make([]*model.StopPoint, 0)
	var byId = make(map[string]*model.StopPoint)

	for _, stop := range g.Stops {

		// FIXME: Nil-checks
		lat, err := strconv.ParseFloat(*stop.Lat, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error parsing shape.PtLat: %v, line: %v", *stop.Lat, stop.LineNumber))
		}
		lon, err := strconv.ParseFloat(*stop.Lon, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error parsing shape.PtLon: %v, line: %v", *stop.Lon, stop.LineNumber))
		}

		lat2 := math.Round(lat*100000) / 100000
		lon2 := math.Round(lon*100000) / 100000

		s := model.StopPoint{
			Name:       *stop.Name,
			ShortName:  *stop.Code,
			Latitude:   lat2,
			Longitude:  lon2,
			TariffZone: *stop.ZoneId,
		}

		if !ggtfs.StringIsNilOrEmpty(stop.Extensions.MunicipalityId) {
			m, err := municipalities.GetOne(*stop.Extensions.MunicipalityId)
			if err != nil {
				warnings = append(warnings, errors.New(fmt.Sprintf("stop-point (%v): municipality information not found, ignoring the stop-point", stop.Id)))
				continue
			}
			s.Municipality = m
		}

		all = append(all, &s)
		byId[*stop.Id] = &s
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].ShortName < all[y].ShortName
	})

	return StopPoints{
		all:  all,
		byId: byId,
	}
}
