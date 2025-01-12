package repository

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"sort"
)

type JourneysLineDataStore struct {
	All  []*model.Line
	ById map[string]*model.Line
}

func newLineDataStore(routes []*ggtfs.Route) *JourneysLineDataStore {
	var all = make([]*model.Line, 0)
	var byId = make(map[string]*model.Line)

	for _, r := range routes {
		id := r.Id

		if id == nil || len(*id) == 0 {
			fmt.Println(fmt.Sprintf("skipping malformed line from routes.txt on line: %v", r.LineNumber))
			continue
		}

		var shortName, longName string

		ln := r.LongName
		if ln == nil || len(*ln) == 0 {
			longName = ""
		} else {
			longName = *ln
		}

		sn := r.ShortName
		if sn == nil || len(*sn) == 0 {
			shortName = ""
		} else {
			shortName = *sn
		}

		l := model.Line{
			Name:        shortName,
			Description: longName,
		}

		all = append(all, &l)
		byId[*id] = &l
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Name < all[y].Name
	})

	return &JourneysLineDataStore{
		All:  all,
		ById: byId,
	}
}
