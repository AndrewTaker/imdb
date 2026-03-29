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

func seedCollections(d *mongo.Database, names []string) error {
	for _, name := range names {
		if err := createCollection(d, name); err != nil {
			return fmt.Errorf("failed to create collection %s", name)
		}
	}

	return nil
}

func ensureIndexes(d *mongo.Database) error {
	userIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := d.Collection("users").Indexes().CreateOne(context.Background(), userIndex)

	genresIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "genres", Value: 1}},
		Options: options.Index().SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2,
		}),
	}
	_, err = d.Collection("movies").Indexes().CreateOne(context.Background(), genresIndex)

	return err
}
