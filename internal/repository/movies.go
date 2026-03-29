package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Movie struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Genres        []string           `bson:"genres"`
	Year          int                `bson:"year"`
	AverageRating float64            `bson:"average_rating"`
	VoteCount     int                `bson:"vote_count"`
}

type Rating struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Score   int                `bson:"score"` // 1-10
}

type MoviesRepository struct {
	c *mongo.Collection
}

func NewMoviesRepository(c *mongo.Collection) *MoviesRepository {
	return &MoviesRepository{c: c}
}

func (r *MoviesRepository) Create(ctx context.Context, title string, genres []string, year int) error {
	if _, err := r.c.InsertOne(ctx, Movie{
		Title:  title,
		Genres: genres,
		Year:   year,
	}); err != nil {
		return err
	}

	return nil
}

func (r *MoviesRepository) GetAll(ctx context.Context, pag PaginationOptions, sort []SortOptions, filter []FilterOptions) ([]Movie, error) {
	opts := CreateQueryOptions(pag, sort, filter)
	filters := CreateFilterOptions(filter)

	rows, err := r.c.Find(ctx, filters, opts)
	if err != nil {
		return nil, err
	}
	defer rows.Close(ctx)

	movies := make([]Movie, 0)
	for rows.Next(ctx) {
		var movie Movie
		if err := rows.Decode(&movie); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MoviesRepository) GetByID(ctx context.Context, id string) (*Movie, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var movie Movie
	err = r.c.FindOne(ctx, bson.M{"_id": objID}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MoviesRepository) AlreadyExists(ctx context.Context, title string, year int) bool {
	filter := bson.D{{"title", title}, {"year", year}}
	result, err := r.c.CountDocuments(ctx, filter, nil)
	if err != nil || result > 0 {
		return true
	}

	return false
}
