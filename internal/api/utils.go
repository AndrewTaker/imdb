package api

import (
	"net/http"
	"strconv"
)

// helper: extract int from query params
// fallback to default value in case of error
func getIntQuery(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

func getStringQuery(r *http.Request, key string, defaultValue string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func getHostWithUri(r *http.Request) string {
	return r.Host + r.RequestURI
}
