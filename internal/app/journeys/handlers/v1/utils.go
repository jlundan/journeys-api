package v1

import "net/http"

func getQueryParameters(r *http.Request) map[string]string {
	result := make(map[string]string)
	for k, v := range r.URL.Query() {
		if k != "exclude-fields" {
			result[k] = v[0]
		}
	}
	return result
}

func getExcludeFieldsQueryParameter(r *http.Request) string {
	if r != nil && r.URL != nil && r.URL.Query() != nil {
		return r.URL.Query().Get("exclude-fields")
	}
	return ""
}
