package api

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MovieResponse struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Genres        []string `json:"genres"`
	Year          int      `json:"year"`
	AverageRating float64  `json:"average_rating"`
	VoteCount     int      `json:"vote_count"`
}

type CreateMovieRequest struct {
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
	Year   int      `json:"year"`
}

type UpdateMovieRequest struct {
	Title  *string   `json:"title" bson:"title,omitempty"`
	Genres *[]string `json:"genres" bson:"genres,omitempty"`
	Year   *int      `json:"year" bson:"year,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application-json")
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(&response)
}
