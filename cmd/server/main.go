package main

import (
	"context"
	"fmt"
	"imdb/internal/database"
	"imdb/internal/repository"
	"imdb/internal/service"
	"log"
	"log/slog"
)

func main() {
	dsn := "mongodb://admin:admin@localhost:27017"

	client, err := database.NewMongoClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.Default()

	err = client.Database("imdb").CreateCollection(context.Background(), "users")
	if err != nil {
		log.Fatal(err)
	}
	ur := repository.NewUserRepository(client.Database("imdb").Collection("users"))
	us := service.NewUserService(ur, logger)
	err = us.Create(context.Background(), "admin@admin.com", "admin")
	if err != nil {
		log.Fatal(err)
	}

	users, err := us.GetAll(context.Background(), repository.PaginationOptions{Limit: 0, Offset: 0}, []repository.UserSortOptions{
		{SortBy: "email", Order: repository.ASC},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

}
