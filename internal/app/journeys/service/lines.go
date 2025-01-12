package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
)

func (ds DataService) SearchLines(params map[string]string) []*model.Line {
	result := make([]*model.Line, 0)

	for _, line := range ds.DataStore.Lines.All {
		if lineMatchesConditions(line, params) {
			result = append(result, line)
		}
	}

	return result
}

func (ds DataService) GetOneLineById(id string) (*model.Line, error) {
	if l, ok := ds.DataStore.Lines.ById[id]; ok {
		return l, nil
	}
	return nil, model.ErrNoSuchElement
}

func lineMatchesConditions(line *model.Line, conditions map[string]string) bool {
	if line == nil {
		return false
	}
	if conditions == nil {
		return true
	}

	for k, v := range conditions {
		switch k {
		case "name":
			if line.Name != v {
				return false
			}
		case "description":
			if !utils.StrContains(line.Description, v) {
				return false
			}
		}
	}

	return true
}
