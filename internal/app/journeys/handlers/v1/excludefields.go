package v1

import (
	"encoding/json"
	"strings"
)

func removeExcludedFields(bodyElements []map[string]any, propertyPaths string) []map[string]any {
	if len(bodyElements) == 0 {
		return bodyElements
	}

	var filteredBodyElements []map[string]any
	for _, be := range bodyElements {
		filteredBodyElements = append(filteredBodyElements, filterMap(be, propertyPaths))
	}
	return filteredBodyElements
}

func convertToStringAnyMap(obj any) (map[string]any, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	objectData := make(map[string]any)

	err = json.Unmarshal(jsonStr, &objectData)
	if err != nil {
		return nil, err
	}

	return objectData, nil
}

func filterMap(obj map[string]any, propertyPaths string) map[string]any {
	ppArr := strings.Split(propertyPaths, ",")
	for _, pp := range ppArr {
		deleteProperty(obj, pp)
	}

	return obj
}

func deleteProperty(obj any, propertyPath string) {
	// propertyPath can be for example: "municipality.name" (stop-points?exclude-fields=municipality.name)
	// pathFragments will then be ["municipality", "name"]
	pathFragments := strings.Split(propertyPath, ".")
	if objectAsMap, isMap := obj.(map[string]interface{}); isMap {
		for i, f := range pathFragments {
			// on the first loop, [i,f] = [0, "municipality"]
			// on the second loop, [i,f] = [1, "name"]
			isLastFragment := i == len(pathFragments)-1
			if _, hasKey := objectAsMap[f]; !hasKey { // key not found - do not delete anything
				return
			} else if isLastFragment {
				// This is the last fragment given. We should delete the key. In our example, this deletes
				// the municipality's id, but the URL could have been stop-points?exclude-fields=municipality,
				// which would then delete the entire municipality from the response.
				delete(objectAsMap, f)
			} else if arr, isArray := objectAsMap[f].([]interface{}); isArray {
				// If the field is an array, for example journeys/<journey-id>?exclude-fields=calls.arrivalTime (and
				// it is not the last item in the pathFragments), then we want to delete the arrivalTime from all
				// calls in the array. Call deleteProperty recursively for all expectedEntities in the array.
				for _, a := range arr {
					deleteProperty(a, strings.Join(pathFragments[i+1:], "."))
				}
			} else if _, isMap2 := objectAsMap[f].(map[string]interface{}); isMap2 {
				// Finally, if the fragment was not the last one, or wasn't an array, we "drill down" to the next level
				// (in our example, we would go from "municipality" to "name"). The for loop will move on to the next
				// path fragment, which will be checked against the new objectAsMap. For example if the obj was initially
				// {
				//		"location": "61.49754,23.76152",
				//		"municipality": {
				//			"name": "Tampere",
				//			"shortName": "837",
				//			"url": "https://data.itsfactory.fi/journeys/api/1/municipalities/837"
				//		},
				//		"name": "Keskustori H",
				//		"shortName": "0001",
				//		"tariffZone": "A",
				//		"url": "https://data.itsfactory.fi/journeys/api/1/stop-points/0001"
				// }
				//
				// the new object will be:
				//
				// {
				//		"name": "Tampere",
				//		"shortName": "837",
				//		"url": "https://data.itsfactory.fi/journeys/api/1/municipalities/837"
				// }
				// since the f was "municipality"
				objectAsMap = objectAsMap[f].(map[string]interface{})
			}
		}
	}
}
