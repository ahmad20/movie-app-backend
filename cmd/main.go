package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"movie-app-go/configs"
	"movie-app-go/entities"
	"movie-app-go/modules/auth"
	"movie-app-go/modules/movie"
	"movie-app-go/modules/user"
	"movie-app-go/repositories"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var (
		movies  []entities.Movie
		users   []entities.User
		tickets []entities.Ticket
	)

	// Load Config
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Load Movies
	response, err := http.Get(config.Data.Movies)
	if err != nil {
		log.Println("Error fetching data:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &movies)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}

	// Set Seats
	for i := range movies {
		movies[i].Seats = repositories.GenerateSeats()
	}

	// Set Router
	router := gin.Default()

	// Set CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Cors.AllowedOrigins
	corsConfig.AllowMethods = config.Cors.AllowedMethods
	router.Use(cors.New(corsConfig))

	authService := auth.NewService(config.JWT.SecretKey, config.JWT.ExpiresIn)
	middleware := auth.AuthMiddleware(authService)

	ticketRepo := repositories.NewTicketRepository(tickets)

	movieRepo := repositories.NewMovieRepository(movies)
	movieUseCase := movie.NewUseCase(movieRepo)
	movieHandler := movie.NewHandler(movieUseCase)
	movie.SetupRouter(router, movieHandler)

	userRepo := repositories.NewUserRepository(users)
	userUseCase := user.NewUseCase(userRepo, movieRepo, ticketRepo)
	userHandler := user.NewHandler(userUseCase, authService)
	user.SetupRouter(router, userHandler, middleware)

	router.Run()
}
