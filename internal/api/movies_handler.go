package api

import (
	"encoding/json"
	"imdb/internal/repository"
	"imdb/internal/security"
	"imdb/internal/service"
	"log/slog"
	"net/http"
)

type MoviesHandler struct {
	s     *service.MoviesService
	l     *slog.Logger
	token *security.TokenService
}

func NewMoviesHandler(s *service.MoviesService, l *slog.Logger, token *security.TokenService) *MoviesHandler {
	return &MoviesHandler{s: s, l: l, token: token}
}

func (h *MoviesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("MoviesHandler.Create", "decoding error", err)
		ErrorResponse(w, http.StatusBadRequest, ErrBadPayload.Error())
		return
	}

	err := h.s.Create(r.Context(), payload.Title, payload.Genres, payload.Year)
	if err != nil {
		h.l.Error("MoviesHandler.Create", "db error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MoviesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var payload UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.l.Error("MoviesHandler.Update", "decoding error", err)
		ErrorResponse(w, http.StatusBadRequest, ErrBadPayload.Error())
		return
	}

	err := h.s.PartialUpdate(r.Context(), id, payload.Title, payload.Year, payload.Genres)
	if err != nil {
		h.l.Error("MoviesHandler.Update", "db error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MoviesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var filters []repository.FilterOptions
	var sorts []repository.SortOptions

	limit := getIntQuery(r, "limit", 10)
	offset := getIntQuery(r, "offset", 0)

	if limit > 25 {
		http.Error(w, "limit capacity is 25", http.StatusBadRequest)
		return
	}

	// filter
	fYear := getIntQuery(r, "f_year", 0)
	if fYear != 0 {
		filters = append(filters, repository.FilterOptions{FilterBy: "year", Value: fYear})
	}
	fGenre := getStringQuery(r, "f_genre", "")
	if fGenre != "" {
		filters = append(filters, repository.FilterOptions{FilterBy: "genres", Value: fGenre})
	}

	// sort
	sorts = append(sorts, ParseSortQuery(r.URL.Query().Get("sort"))...)

	movies, err := h.s.GetAll(
		r.Context(),
		repository.PaginationOptions{Limit: limit, Offset: offset},
		sorts,
		filters,
	)
	if err != nil {
		h.l.Error("MoviesHandler.GetAll", "error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
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
			VoteCount:     m.VoteCount,
		}
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(moviesResponse)
	if err != nil {
		h.l.Error("MoviesHandler.GetAll", "error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}
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
		ErrorResponse(w, http.StatusNotFound, "movie not found")
		return
	}

	var movieResponse MovieResponse
	movieResponse.Title = movie.Title
	movieResponse.Genres = movie.Genres
	movieResponse.ID = movie.ID.Hex()
	movieResponse.Year = movie.Year
	movieResponse.AverageRating = movie.AverageRating
	movieResponse.VoteCount = movie.VoteCount

	err = json.NewEncoder(w).Encode(movieResponse)
	if err != nil {
		h.l.Error("MoviesHandler.GetByID", "error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}
}

func (h *MoviesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrBadPath.Error())
		return
	}

	err := h.s.DeleteByID(r.Context(), id)
	if err != nil {
		h.l.Error("MoviesHandler.Delete", "error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MoviesHandler) Rate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		ErrorResponse(w, http.StatusBadRequest, ErrBadPath.Error())
		return
	}

	score := getIntParam(r, "score")
	if !validRatingScore(score) {
		ErrorResponse(w, http.StatusBadRequest, "invalid rating value, must be >= 1 and <= 10")
		return
	}

	userID, err := h.token.GetUserID(r)
	if userID == "" || err != nil {
		ErrorResponse(w, http.StatusUnauthorized, ErrUnauthorized.Error())
		return
	}

	err = h.s.Rate(r.Context(), userID, id, score)
	if err != nil {
		h.l.Error("MoviesHandler.Rate", "error", err)
		ErrorResponse(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}
}
