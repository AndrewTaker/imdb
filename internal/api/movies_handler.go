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
	limit := getIntQuery(r, "limit", 10)
	offset := getIntQuery(r, "offset", 0)
	h.l.Info("pag values", "offset", offset, "limit", limit)

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

	moviesResponse := make([]MoviesResponse, len(movies))
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
