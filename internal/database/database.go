package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(dsn string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(dsn).
		SetConnectTimeout(10 * time.Second).
		SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongo ping error: %w", err)
	}

	collections := []string{"users", "movies"}
	err = seedCollections(client.Database("imdb"), collections)
	if err != nil {
		return nil, err
	}

	err = ensureIndexes(client.Database("imdb"))
	if err != nil {
		return nil, err
	}

	return client, nil
}
