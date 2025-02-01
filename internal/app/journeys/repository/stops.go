package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"math"
	"sort"
	"strconv"
)

func newStopPointsRepository(stops []*ggtfs.Stop, municipalityDataStore *JourneysMunicipalitiesRepository) *JourneysStopPointsRepository {
	var all = make([]*model.StopPoint, 0)
	var byId = make(map[string]*model.StopPoint)

	for _, stop := range stops {

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
			if m, ok := municipalityDataStore.ById[*stop.Extensions.MunicipalityId]; ok {
				s.Municipality = m
			} else {
				fmt.Println(fmt.Sprintf("stop-point (%v): municipality information not found, ignoring the stop-point", stop.Id))
				continue
			}
		}

		all = append(all, &s)
		byId[*stop.Id] = &s
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].ShortName < all[y].ShortName
	})

	return &JourneysStopPointsRepository{
		All:  all,
		ById: byId,
	}
}

type JourneysStopPointsRepository struct {
	All  []*model.StopPoint
	ById map[string]*model.StopPoint
}
