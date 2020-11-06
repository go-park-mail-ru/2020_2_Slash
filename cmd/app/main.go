package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"

	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"

	sessionHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery"
	sessionRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/session/repository"
	sessionUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/session/usecases"

	userHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery"
	userRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/user/usecases"

	genreHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/delivery"
	genreRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/repository"
	genreUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/usecases"

	countryHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/country/delivery"
	countryRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/country/repository"
	countryUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/country/usecases"

	actorHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/delivery"
	actorRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/repository"
	actorUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/usecases"

	directorHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/director/delivery"
	directorRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/director/repository"
	directorUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/director/usecases"

	contentRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/content/repository"
	contentUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/content/usecases"

	movieHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/delivery"
	movieRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/repository"
	movieUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/usecases"
	ratingHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/delivery"
	ratingRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/repository"
	ratingUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/usecases"
)

func main() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Locale storage
	avatarsPath := config.GetAvatarsPath()
	helpers.InitStorage(avatarsPath)

	postersPath := config.GetPostersPath()
	helpers.InitStorage(postersPath)

	videosPath := config.GetVideosPath()
	helpers.InitStorage(videosPath)

	// Database
	dbConnection, err := sql.Open("postgres", config.GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	// Repository
	sessRepo := sessionRepo.NewSessionPgRepository(dbConnection)
	userRepo := userRepo.NewUserPgRepository(dbConnection)
	genreRepo := genreRepo.NewGenrePgRepository(dbConnection)
	countryRepo := countryRepo.NewCountryPgRepository(dbConnection)
	actorRepo := actorRepo.NewActorPgRepository(dbConnection)
	directorRepo := directorRepo.NewDirectorPgRepository(dbConnection)
	contentRepo := contentRepo.NewContentPgRepository(dbConnection)
	movieRepo := movieRepo.NewMoviePgRepository(dbConnection)
	ratingRepo := ratingRepo.NewRatingPgRepository(dbConnection)

	// Usecases
	sessUcase := sessionUsecase.NewSessionUsecase(sessRepo)
	userUcase := userUsecase.NewUserUsecase(userRepo)
	genreUcase := genreUsecase.NewGenreUsecase(genreRepo)
	countryUcase := countryUsecase.NewCountryUsecase(countryRepo)
	actorUcase := actorUsecase.NewActorUseCase(actorRepo)
	directorUcase := directorUsecase.NewDirectorUseCase(directorRepo)
	contentUcase := contentUsecase.NewContentUsecase(contentRepo, countryUcase, genreUcase, actorUcase, directorUcase)
	movieUcase := movieUsecase.NewMovieUsecase(movieRepo, contentUcase)
	ratingUcase := ratingUsecase.NewRatingUseCase(ratingRepo, contentUcase)

	// Middleware
	e := echo.New()
	mw := mwares.NewMiddlewareManager(sessUcase)
	e.Use(mw.PanicRecovering, mw.AccessLog, mw.CORS)

	e.Static("/avatars", avatarsPath)
	e.Static("/images", postersPath)
	e.Static("/videos", videosPath)

	// Delivery
	sessionHandler := sessionHandler.NewSessionHandler(sessUcase, userUcase)
	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	genreHandler := genreHandler.NewGenreHandler(genreUcase)
	countryHandler := countryHandler.NewCountryHandler(countryUcase)
	actorHandler := actorHandler.NewActorHandler(actorUcase)
	directorHandler := directorHandler.NewDirectorHandler(directorUcase)
	movieHandler := movieHandler.NewMovieHandler(movieUcase, contentUcase, countryUcase, genreUcase, actorUcase, directorUcase)
	ratingHandler := ratingHandler.NewRatingHandler(ratingUcase)

	userHandler.Configure(e, mw)
	sessionHandler.Configure(e, mw)
	genreHandler.Configure(e, mw)
	countryHandler.Configure(e, mw)
	actorHandler.Configure(e, mw)
	directorHandler.Configure(e, mw)
	movieHandler.Configure(e, mw)
	ratingHandler.Configure(e, mw)

	log.Fatal(e.Start(config.GetServerConnString()))
}
