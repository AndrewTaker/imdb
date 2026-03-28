package service

import (
	"context"
	"imdb/internal/repository"
)

type MoviesService struct {
	r *repository.MoviesRepository
}

func NewMoviesService(r *repository.MoviesRepository) *MoviesService {
	return &MoviesService{r: r}
}

func (s *MoviesService) GetAll(
	ctx context.Context,
	pag repository.PaginationOptions,
	sort []repository.SortOptions,
	filter []repository.FilterOptions,
) ([]repository.Movie, error) {
	return s.r.GetAll(ctx, pag, sort, filter)
}

func (s *MoviesService) GetByID(ctx context.Context, id string) (*repository.Movie, error) {
	return s.r.GetByID(ctx, id)
}

func (s *MoviesService) Create(ctx context.Context, title string, genres []string, year int) error {
	return s.r.Create(ctx, title, genres, year)
}
