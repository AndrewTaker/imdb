package service

import (
	"context"
	"fmt"
	"imdb/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MoviesService struct {
	r      *repository.MoviesRepository
	rating *repository.RatingRepository
}

func NewMoviesService(r *repository.MoviesRepository, rating *repository.RatingRepository) *MoviesService {
	return &MoviesService{r: r, rating: rating}
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	return s.r.GetByID(ctx, objID)
}

func (s *MoviesService) Create(ctx context.Context, title string, genres []string, year int) error {
	// not sure if i should do that performance wise
	if exists := s.r.AlreadyExists(ctx, title, year); exists {
		return fmt.Errorf("movie with title %s of year %d already exists", title, year)
	}

	return s.r.Create(ctx, title, genres, year)
}

func (s *MoviesService) PartialUpdate(ctx context.Context, id string, title *string, year *int, genres *[]string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.r.PartialUpdate(ctx, objID, title, year, genres)
}

func (s *MoviesService) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	if err := s.r.Delete(ctx, objID); err != nil {
		return err
	}

	return nil
}

// this should probably be a transaction
// https://www.mongodb.com/docs/drivers/go/current/crud/transactions/
func (s *MoviesService) Rate(ctx context.Context, userID, movieID string, score int) error {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	mID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return fmt.Errorf("invalid movie id format: %w", err)
	}

	if err := s.rating.UpsertRating(ctx, mID, uID, score); err != nil {
		return fmt.Errorf("failed to save rating: %w", err)
	}

	stats, err := s.rating.CalculateRatingStats(ctx, mID)
	if err != nil {
		return fmt.Errorf("failed to calculate new stats: %w", err)
	}

	if err := s.r.UpdateStats(ctx, mID, stats.AverageRating, stats.VoteCount); err != nil {
		return fmt.Errorf("failed to sync movie stats: %w", err)
	}

	return nil
}
