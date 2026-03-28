package api

import (
	"encoding/json"
	"imdb/internal/repository"
	"imdb/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type MoviesHandler struct {
	s *service.MoviesService
	l *slog.Logger
}

func NewMoviesHandler(s *service.MoviesService, l *slog.Logger) *MoviesHandler {
	return &MoviesHandler{s: s, l: l}
}

func (h *MoviesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.l.Info("/movies", "method", r.Method)
	var limit int = 10
	var offset int = 0
	var err error

	limitQuery := r.URL.Query().Get("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			h.l.Error("err casting atoi for limit %s", limitQuery)
			http.Error(w, "bad limit value", http.StatusBadRequest)
			return
		}
	}

	offsetQuery := r.URL.Query().Get("limit")
	if offsetQuery != "" {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			h.l.Error("err casting atoi for offset %s", offsetQuery)
			http.Error(w, "bad offset value", http.StatusBadRequest)
			return
		}
	}

	movies, err := h.s.GetAll(
		r.Context(),
		repository.PaginationOptions{Limit: limit, Offset: offset},
		nil,
		nil,
	)
	if err != nil {
		h.l.Error("err getting movies", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")

	var moviesResponse []MoviesResponse
	for _, m := range movies {
		moviesResponse = append(moviesResponse, MoviesResponse{
			ID:     m.ID.Hex(),
			Title:  m.Title,
			Year:   m.Year,
			Genres: m.Genres,
			Rating: m.Rating,
		})
	}
	if err := json.NewEncoder(w).Encode(moviesResponse); err != nil {
		h.l.Error("err encoding movies", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}
}
