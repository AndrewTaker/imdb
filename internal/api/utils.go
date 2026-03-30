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

func getIntParam(r *http.Request, key string) int {
	val := r.PathValue(key)
	if val == "" {
		return -1
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}

	return intVal
}

func getHostWithUri(r *http.Request) string {
	return r.Host + r.RequestURI
}

func validRatingScore(s int) bool {
	if s == -1 {
		return false
	}
	if s < 1 || s > 10 {
		return false
	}

	return true
}
