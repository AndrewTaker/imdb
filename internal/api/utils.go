package api

import (
	"imdb/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

// extract string query as float64
func getFloat64Query(r *http.Request, key string, defaultValue float64) float64 {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// extract string query as integer
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

// extract string query as string
func getStringQuery(r *http.Request, key string, defaultValue string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}

	return val
}

// extract string param as integer
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

// validate movie struct score field
// must be >= 1 && <= 10
func validRatingScore(s int) bool {
	if s == -1 {
		return false
	}
	if s < 1 || s > 10 {
		return false
	}

	return true
}

// parse sort query from url
// example query is <url>/?sort=title:asc,year:desc,average_rating:asc
// default order is asc
func ParseSortQuery(raw string) []repository.SortOptions {
	if raw == "" {
		return nil
	}

	var sorts []repository.SortOptions
	// TODO: gopls says it is better to use SplitSeq
	// but we have at most 3 values
	pairs := strings.Split(raw, ",")

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		fieldName := strings.TrimSpace(parts[0])
		if fieldName == "" {
			continue
		}

		s := repository.SortOptions{
			SortBy: fieldName,
			Order:  repository.ASC,
		}

		if len(parts) > 1 {
			switch strings.ToLower(parts[1]) {
			case "desc":
				s.Order = repository.DESC
			case "asc":
				s.Order = repository.ASC
			}
		}

		sorts = append(sorts, s)
	}

	return sorts
}
