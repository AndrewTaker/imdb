package repository

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RatingRepository struct {
	c *mongo.Collection
}

func NewRatingRepository(c *mongo.Collection) *RatingRepository {
	return &RatingRepository{c: c}
}

func (r *RatingRepository) UpsertRating(ctx context.Context, movieID, userID primitive.ObjectID, score int) error {
	filter := bson.D{{"movie_id", movieID}, {"user_id", userID}}
	update := bson.D{{"$set", bson.D{{"score", score}}}}

	opts := options.Update().SetUpsert(true)

	_, err := r.c.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *RatingRepository) CalculateRatingStats(ctx context.Context, movieID primitive.ObjectID) (*RatingStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"movie_id": movieID}}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"average_rating": bson.M{"$avg": "$score"},
			"vote_count":     bson.M{"$sum": 1},
		}}},
	}
	fmt.Println(pipeline)

	cursor, err := r.c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var s RatingStats
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		log.Println(s)
		return &s, nil
	}

	return &RatingStats{AverageRating: 0, VoteCount: 0}, nil
}
