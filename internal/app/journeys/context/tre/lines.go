package tre

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"sort"
)

type Lines struct {
	all  []*model.Line
	byId map[string]*model.Line
}

func (lines Lines) GetOne(name string) (*model.Line, error) {
	if _, ok := lines.byId[name]; !ok {
		return &model.Line{}, model.ErrNoSuchElement
	}
	return lines.byId[name], nil
}
func (lines Lines) GetAll() []*model.Line {
	return lines.all
}

func buildLines(g GTFSContext) Lines {
	var all = make([]*model.Line, 0)
	var byId = make(map[string]*model.Line)

	for _, v := range g.Routes {
		ln := v.LongName
		sn := v.ShortName

		if !ln.IsValid() || !sn.IsValid() {
			fmt.Println(fmt.Sprintf("skipping malformed line: %v", v.Id))
			continue
		}

		l := model.Line{
			Name:        sn.String(),
			Description: ln.String(),
		}

		all = append(all, &l)
		byId[sn.String()] = &l
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Name < all[y].Name
	})

	return Lines{
		all:  all,
		byId: byId,
	}
}
