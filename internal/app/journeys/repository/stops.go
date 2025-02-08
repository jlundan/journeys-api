package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

func newStopPointsRepository(stops []*ggtfs.Stop, municipalityDataStore *JourneysMunicipalitiesRepository) *JourneysStopPointsRepository {
	var all = make([]*model.StopPoint, 0)
	var byId = make(map[string]*model.StopPoint)

	for _, stop := range stops {
		var lat, lon float64
		if stop.Lat != nil {
			lf, err := strconv.ParseFloat(*stop.Lat, 64)
			if err != nil {
				log.Println(fmt.Sprintf("stop-point (on gtfs line %v): cannot parse lat float value", stop.LineNumber))
			}
			lat = lf
		} else {
			log.Println(fmt.Sprintf("stop-point (on gtfs line %v): lat is missing", stop.LineNumber))
		}

		if stop.Lon != nil {
			lf, err := strconv.ParseFloat(*stop.Lon, 64)
			if err != nil {
				log.Println(fmt.Sprintf("stop-point (on gtfs line %v): cannot parse lon float value", stop.LineNumber))
			}
			lon = lf
		} else {
			log.Println(fmt.Sprintf("stop-point (on gtfs line %v): lon is missing", stop.LineNumber))
		}

		var name, shortName, tariffZone string

		if stop.Name != nil {
			name = strings.TrimSpace(*stop.Name)
		} else {
			log.Println(fmt.Sprintf("stop-point (on gtfs line %v): name is missing", stop.LineNumber))
		}

		if stop.Code != nil {
			shortName = strings.TrimSpace(*stop.Code)
		} else {
			log.Println(fmt.Sprintf("stop-point (on gtfs line %v): shortName is missing", stop.LineNumber))
		}

		if stop.ZoneId != nil {
			tariffZone = strings.TrimSpace(*stop.ZoneId)
		} else {
			log.Println(fmt.Sprintf("stop-point (on gtfs line %v): tariffZone is missing", stop.LineNumber))
		}

		s := model.StopPoint{
			Name:       name,
			ShortName:  shortName,
			Latitude:   math.Round(lat*100000) / 100000,
			Longitude:  math.Round(lon*100000) / 100000,
			TariffZone: tariffZone,
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
