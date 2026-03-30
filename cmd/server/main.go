package main

import (
	"context"
	"imdb/internal/api"
	"imdb/internal/database"
	"imdb/internal/repository"
	"imdb/internal/security"
	"imdb/internal/service"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	dsn := "mongodb://admin:admin@localhost:27017"

	client, err := database.NewMongoClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	logger := slog.Default()

	ts, err := security.NewTokenService("12345678901234567890123456789012")
	if err != nil {
		log.Fatal(err)
	}
	ur := repository.NewUserRepository(client.Database("imdb").Collection("users"))
	us := service.NewUserService(ur, logger)
	uh := api.NewUsersHandler(us, logger, ts)

	mr := repository.NewMoviesRepository(client.Database("imdb").Collection("movies"))
	rr := repository.NewRatingRepository(client.Database("imdb").Collection("ratings"))
	ms := service.NewMoviesService(mr, rr)
	mh := api.NewMoviesHandler(ms, logger, ts)

	creds := map[string]string{
		"admin@admin.com": "admin",
		"user001":         "user001",
		"user002":         "user002",
	}
	for email, password := range creds {
		_ = us.Create(context.Background(), email, password)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/signin", uh.SignIn)
	mux.HandleFunc("POST /auth/signup", uh.SignUp)
	mux.HandleFunc("POST /movies", mh.Create)
	mux.HandleFunc("GET /movies", mh.GetAll)
	mux.HandleFunc("GET /movies/{id}", mh.GetByID)
	mux.HandleFunc("DELETE /movies/{id}", mh.Delete)
	mux.HandleFunc("PATCH /movies/{id}", mh.Update)
	mux.HandleFunc("POST /movies/{id}/rate/{score}", mh.Rate) // auth required

	log.Println("listening on 4444")
	log.Fatal(http.ListenAndServe(":4444", api.LogRequestMiddleware(mux)))

}
