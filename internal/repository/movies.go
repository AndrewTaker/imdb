package repository

import (
	"context"

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

func (r *MoviesRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if _, err := r.c.DeleteOne(ctx, bson.D{{"_id", id}}); err != nil {
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

func (r *MoviesRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Movie, error) {
	var movie Movie
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MoviesRepository) PartialUpdate(ctx context.Context, id primitive.ObjectID, title *string, year *int, genres *[]string) error {
	// if genres are not nil
	// we append to set rather than rewriting values
	// https://www.mongodb.com/docs/manual/reference/operator/update/addToSet/
	set := bson.M{}
	addToSet := bson.M{}

	if title != nil {
		set["title"] = *title
	}
	if year != nil {
		set["year"] = *year
	}

	if genres != nil && len(*genres) > 0 {
		// https://www.mongodb.com/docs/manual/reference/operator/update/each/
		addToSet["genres"] = bson.M{"$each": *genres}
	}

	update := bson.M{}
	if len(set) > 0 {
		update["$set"] = set
	}
	if len(addToSet) > 0 {
		update["$addToSet"] = addToSet
	}

	if len(update) == 0 {
		return nil
	}

	_, err := r.c.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MoviesRepository) AlreadyExists(ctx context.Context, title string, year int) bool {
	filter := bson.D{{"title", title}, {"year", year}}
	result, err := r.c.CountDocuments(ctx, filter, nil)
	if err != nil || result > 0 {
		return true
	}

	return false
}
