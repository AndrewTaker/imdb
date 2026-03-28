package api

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MoviesResponse struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
}
