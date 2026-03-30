package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createCollection(d *mongo.Database, name string) error {
	return d.CreateCollection(context.Background(), name)
}

// function to create collections
// called inside NewDatabase()
func seedCollections(d *mongo.Database, names []string) error {
	for _, name := range names {
		if err := createCollection(d, name); err != nil {
			return fmt.Errorf("failed to create collection %s", name)
		}
	}

	return nil
}

// function to create indexes
// since we only have a few - we can use single function
// called inside NewDatabase()
func ensureIndexes(d *mongo.Database) error {
	// unique constraint for user's email
	userIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := d.Collection("users").Indexes().CreateOne(context.Background(), userIndex)

	// case insensetive search for genres
	genresIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "genres", Value: 1}},
		Options: options.Index().SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	}

	// unique constraint for title + year
	movieIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: 1},
			{Key: "year", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = d.Collection("movies").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		genresIndex,
		movieIndex,
	})

	return err
}
