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

func (h *MoviesHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.l.Info(getHostWithUri(r), "method", r.Method)

	var payload CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("MoviesHandler.Create", "decoding error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.s.Create(r.Context(), payload.Title, payload.Genres, payload.Year)
	if err != nil {
		h.l.Error("MoviesHandler.Create", "db error", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MoviesHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.l.Info(getHostWithUri(r), "method", r.Method)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var payload UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("MoviesHandler.Update", "decoding error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.s.PartialUpdate(r.Context(), id, payload.Title, payload.Year, payload.Genres)
	if err != nil {
		h.l.Error("MoviesHandler.Update", "db error", err)
		http.Error(w, "could not update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
		h.l.Error("MoviesHandler.GetAll", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}

	moviesResponse := make([]MovieResponse, len(movies))
	for i, m := range movies {
		moviesResponse[i] = MovieResponse{
			ID:            m.ID.Hex(),
			Title:         m.Title,
			Year:          m.Year,
			Genres:        m.Genres,
			AverageRating: m.AverageRating,
		}
	}

	err = json.NewEncoder(w).Encode(moviesResponse)
	if err != nil {
		h.l.Error("MoviesHandler.GetAll", "error", err)
		http.Error(w, "could not get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MoviesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	movie, err := h.s.GetByID(r.Context(), id)
	if err != nil {
		h.l.Error("MoviesHandler.GetByID", "error", err)
		http.Error(w, "movie not found", http.StatusNotFound)
		return
	}

	var movieResponse MovieResponse
	movieResponse.Title = movie.Title
	movieResponse.Genres = movie.Genres
	movieResponse.ID = movie.ID.Hex()
	movieResponse.Year = movie.Year
	movieResponse.AverageRating = movie.AverageRating

	err = json.NewEncoder(w).Encode(movieResponse)
	if err != nil {
		h.l.Error("MoviesHandler.GetByID", "error", err)
		http.Error(w, "could not get movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MoviesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.s.DeleteByID(r.Context(), id)
	if err != nil {
		h.l.Error("MoviesHandler.Delete", "error", err)
		http.Error(w, "could not delete movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
