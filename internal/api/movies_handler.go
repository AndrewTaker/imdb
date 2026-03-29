package api

import (
	"encoding/json"
	"imdb/internal/repository"
	"imdb/internal/service"
	"log/slog"
	"net/http"
)

type MoviesHandler struct {
	s *service.MoviesService
	l *slog.Logger
}

func NewMoviesHandler(s *service.MoviesService, l *slog.Logger) *MoviesHandler {
	return &MoviesHandler{s: s, l: l}
}

func (h *MoviesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.l.Info(getHostWithUri(r), "method", r.Method)
	var filters []repository.FilterOptions

	limit := getIntQuery(r, "limit", 10)
	offset := getIntQuery(r, "offset", 0)

	if limit > 25 {
		http.Error(w, "limit capacity is 25", http.StatusBadRequest)
		return
	}

	year := getIntQuery(r, "year", 0)
	if year != 0 {
		filters = append(filters, repository.FilterOptions{FilterBy: "year", Value: year})
	}
	genre := getStringQuery(r, "genre", "")
	if genre != "" {
		filters = append(filters, repository.FilterOptions{FilterBy: "genres", Value: genre})
	}

	movies, err := h.s.GetAll(
		r.Context(),
		repository.PaginationOptions{Limit: limit, Offset: offset},
		nil,
		filters,
	)
	if err != nil {
		h.l.Error("err getting movies", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}

	moviesResponse := make([]MoviesResponse, len(movies))
	for i, m := range movies {
		moviesResponse[i] = MoviesResponse{
			ID:     m.ID.Hex(),
			Title:  m.Title,
			Year:   m.Year,
			Genres: m.Genres,
			Rating: m.Rating,
		}
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(moviesResponse); err != nil {
		h.l.Error("err encoding movies", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}
}
