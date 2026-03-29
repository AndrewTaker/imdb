package api

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
}

type CreateMovieRequest struct {
	Title  string   `json:"string"`
	Genres []string `json:"genres"`
	Year   int      `json:"year"`
}
