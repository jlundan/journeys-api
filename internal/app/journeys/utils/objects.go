package utils

import (
	"encoding/json"
	"strings"
)

var (
	jsonMarshal   func(_ interface{}) ([]byte, error)
	jsonUnmarshal func(data []byte, v interface{}) error
)

func init() {
	jsonMarshal = json.Marshal
	jsonUnmarshal = json.Unmarshal
}

func ToInterfaceMapViaJSON(obj interface{}) (map[string]interface{}, error) {
	jsonStr, err := jsonMarshal(obj)
	if err != nil {
		return nil, err
	}
	objectData := make(map[string]interface{})
	err = jsonUnmarshal(jsonStr, &objectData)
	if err != nil {
		return nil, err
	}

	return objectData, nil
}

func FilterObject(r interface{}, properties string) (interface{}, error) {
	im, err := ToInterfaceMapViaJSON(r)
	if err != nil {
		return nil, err
	}

	props := strings.Split(properties, ",")
	for _, prop := range props {
		deleteProperty(im, prop)
	}

	return im, nil
}

func deleteProperty(raw interface{}, path string) {
	fragments := strings.Split(path, ".")
	if objectData, isMap := raw.(map[string]interface{}); isMap {
		for i, f := range fragments {
			isLastFragment := i == len(fragments)-1
			if _, hasKey := objectData[f]; !hasKey { // key not found - do not delete anything
				return
			} else if isLastFragment { // is last fragment - delete regardless what it is
				delete(objectData, f)
			} else if arr, isArray := objectData[f].([]interface{}); isArray { // array - recursive call for all items
				for _, a := range arr {
					deleteProperty(a, strings.Join(fragments[i+1:], "."))
				}
			} else if _, isMap := objectData[f].(map[string]interface{}); isMap { // map - drill down to the sub element
				objectData = objectData[f].(map[string]interface{})
			}
		}
	}
}
