package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
)

type DataService struct {
	Context model.Context
}

func (ds DataService) SearchLines(params map[string]string) []*model.Line {
	return make([]*model.Line, 0)
}

func (ds DataService) GetOneLineByName(name string) *model.Line {
	return nil
}

//func GetAllLines(w http.ResponseWriter, r *http.Request, context model.Context) {
//	responseItems := make([]interface{}, 0)
//
//	for _, line := range context.Lines().GetAll() {
//		if lineMatchesConditions(line, getDefaultConditions(r)) {
//			responseItems = append(responseItems, convertLine(line))
//		}
//	}
//
//	sendResponse(responseItems, nil, r, w)
//}
//
//func handleGetOneLine(w http.ResponseWriter, r *http.Request, context model.Context, name string) {
//	line, err := context.Lines().GetOne(name)
//
//	if err == nil {
//		sendResponse([]interface{}{convertLine(line)}, nil, r, w)
//	} else {
//		sendResponse(nil, err, r, w)
//	}
//}
//
//func convertLine(line *model.Line) Line {
//	return Line{
//		Url:         fmt.Sprintf("%v%v/%v", os.Getenv("JOURNEYS_BASE_URL"), linePrefix, line.Name),
//		Name:        line.Name,
//		Description: line.Description,
//	}
//}
//
//func lineMatchesConditions(line *model.Line, conditions map[string]string) bool {
//	if line == nil {
//		return false
//	}
//	if conditions == nil {
//		return true
//	}
//
//	for k, v := range conditions {
//		switch k {
//		case "name":
//			if line.Name != v {
//				return false
//			}
//		case "description":
//			if !utils.StrContains(line.Description, v) {
//				return false
//			}
//		}
//	}
//
//	return true
//}
