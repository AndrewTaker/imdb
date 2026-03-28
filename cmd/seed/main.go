package main

import (
	"context"
	"encoding/json"
	"fmt"
	"imdb/internal/database"
	"imdb/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Movie struct {
	Title string
	Year  string
	Genre string
}

func main() {
	token := os.Getenv("OMDB_TOKEN")
	if token == "" {
		log.Fatal("no token")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	dsn := "mongodb://admin:admin@localhost:27017"
	mongo, err := database.NewMongoClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
	mr := repository.NewMoviesRepository(mongo.Database("imdb").Collection("movies"))

	for i, movie := range movies {
		log.Printf("%d:%s\n", i, movie)
		req, err := http.NewRequest(http.MethodGet, "http://www.omdbapi.com/demo.aspx/", nil)
		if err != nil {
			log.Fatalf("failed to create request: %v", err)
		}

		q := req.URL.Query()
		q.Add("t", movie)
		// q.Add("y", "2025")
		q.Add("plot", "short")
		q.Add("token", token)
		req.URL.RawQuery = q.Encode()

		req.Header.Set("User-Agent", "Go-OMDB-Client/1.0")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected status code: %s %d\n", movie, resp.StatusCode)
			continue
		}

		var m Movie
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			log.Printf("failed to decode: %s\n", movie)
			continue
		}

		mYear, err := strconv.Atoi(m.Year)
		if err != nil {
			log.Printf("failed to convert a -> i %s %s\n", movie, err.Error())
		}
		mGenres := strings.Split(m.Genre, ",")

		err = mr.Create(context.Background(), m.Title, mGenres, mYear)
		if err != nil {
			log.Printf("could not save %s %s\n", movie, err.Error())
		}
	}

	result, _ := mr.GetAll(context.Background(), repository.PaginationOptions{Limit: 0, Offset: 0}, nil, nil)
	fmt.Println(result)
}

var movies = []string{
	"The Shawshank Redemption",
	"The Godfather",
	"The Dark Knight",
	"Pulp Fiction",
	"Schindler's List",
	"Inception",
	"Fight Club",
	"Forrest Gump",
	"The Matrix",
	"Goodfellas",
	"Seven Samurai",
	"Se7en",
	"The Silence of the Lambs",
	"City of God",
	"It's a Wonderful Life",
	"Life Is Beautiful",
	"The Usual Suspects",
	"Léon: The Professional",
	"Spirited Away",
	"Saving Private Ryan",
	"Interstellar",
	"The Green Mile",
	"Parasite",
	"The Prestige",
	"The Lion King",
	"The Departed",
	"The Pianist",
	"Gladiator",
	"Whiplash",
	"The Intouchables",
	"The Godfather Part II",
	"Back to the Future",
	"Psycho",
	"Casablanca",
	"Modern Times",
	"Rear Window",
	"Raiders of the Lost Ark",
	"Apocalypse Now",
	"Alien",
	"Memento",
	"Django Unchained",
	"The Great Dictator",
	"Sunset Boulevard",
	"Paths of Glory",
	"The Shining",
	"WALL-E",
	"American History X",
	"Princess Mononoke",
	"Oldboy",
	"Witness for the Prosecution",
	"Once Upon a Time in the West",
	"Das Boot",
	"Citizen Kane",
	"Vertigo",
	"North by Northwest",
	"Reservoir Dogs",
	"Braveheart",
	"Amadeus",
	"Requiem for a Dream",
	"2001: A Space Odyssey",
	"Lawrence of Arabia",
	"Eternal Sunshine of the Spotless Mind",
	"A Clockwork Orange",
	"Taxi Driver",
	"Full Metal Jacket",
	"Double Indemnity",
	"Toy Story",
	"The Sting",
	"The Apartment",
	"Metropolis",
	"To Kill a Mockingbird",
	"Up",
	"Heat",
	"L.A. Confidential",
	"Die Hard",
	"Snatch",
	"Indiana Jones and the Last Crusade",
	"1917",
	"The Kid",
	"Blade Runner",
	"Unforgiven",
	"The Hunt",
	"The Wolf of Wall Street",
	"No Country for Old Men",
	"Jurrasic Park",
	"Truman Show",
	"Gone with the Wind",
	"Gran Torino",
	"The Bridge on the River Kwai",
	"The Third Man",
	"Blade Runner 2049",
	"Mad Max: Fury Road",
	"Jaws",
	"Arrival",
	"The Grand Budapest Hotel",
	"Portrait of a Lady on Fire",
	"Coco",
	"The Thing",
	"Logan",
	"The Truman Show",
}
