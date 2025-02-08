package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"log"
	"sort"
	"strings"
)

func newLinesRepository(routes []*ggtfs.Route) *JourneysLinesRepository {
	var all = make([]*model.Line, 0)
	var byId = make(map[string]*model.Line)

	for i, r := range routes {
		if r == nil {
			fmt.Println(fmt.Sprintf("Nil route detected, number %v in the routes array, newLinesRepository function", i))
			continue
		}

		if r.Id == nil {
			fmt.Println(fmt.Sprintf("Line.Id is missing, GTFS route: %v", r.LineNumber))
			continue
		}

		id := strings.TrimSpace(*r.Id)

		var shortName, longName string

		if r.LongName == nil {
			log.Println(fmt.Sprintf("line (on gtfs line %v): LongName is missing", r.LineNumber))
		} else {
			longName = strings.TrimSpace(*r.LongName)
		}

		if r.ShortName == nil {
			log.Println(fmt.Sprintf("line (on gtfs line %v): ShortName is missing", r.LineNumber))
		} else {
			shortName = strings.TrimSpace(*r.ShortName)
		}

		l := model.Line{
			Name:        shortName,
			Description: longName,
		}

		all = append(all, &l)
		byId[id] = &l
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Name < all[y].Name
	})

	return &JourneysLinesRepository{
		All:  all,
		ById: byId,
	}
}

type JourneysLinesRepository struct {
	All  []*model.Line
	ById map[string]*model.Line
}
