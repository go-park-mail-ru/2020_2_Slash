package main

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"google.golang.org/grpc"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares/monitoring"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"

	sessionHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery"
	sessionGRPC "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery/grpc"
	sessionUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/session/usecases"

	userHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery"
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

	contentHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/content/delivery"
	contentRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/content/repository"
	contentUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/content/usecases"

	movieHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/delivery"
	movieRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/repository"
	movieUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/usecases"

	tvshowHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/delivery"
	tvshowRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/repository"
	tvshowUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/usecases"

	ratingHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/delivery"
	ratingRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/repository"
	ratingUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/rating/usecases"

	favouriteHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/delivery"
	favouriteRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/repository"
	favouriteUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/usecases"

	seasonHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/season/delivery"
	seasonRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/season/repository"
	seasonUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/season/usecases"

	episodeHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/episode/delivery"
	episodeRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/episode/repository"
	episodeUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/episode/usecases"

	searchHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/search/delivery"
	searchUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/search/usecases"
)

func main() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	logger.InitLogger(config.GetLoggerDir(), config.GetLogLevel())

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
	genreRepo := genreRepo.NewGenrePgRepository(dbConnection)
	countryRepo := countryRepo.NewCountryPgRepository(dbConnection)
	actorRepo := actorRepo.NewActorPgRepository(dbConnection)
	directorRepo := directorRepo.NewDirectorPgRepository(dbConnection)
	contentRepo := contentRepo.NewContentPgRepository(dbConnection)
	movieRepo := movieRepo.NewMoviePgRepository(dbConnection)
	tvshowRepo := tvshowRepo.NewTVShowPgRepository(dbConnection)
	ratingRepo := ratingRepo.NewRatingPgRepository(dbConnection)
	favouriteRepo := favouriteRepo.NewFavouritePgRepository(dbConnection)
	seasonRepo := seasonRepo.NewSeasonPgRepository(dbConnection)
	episodeRepo := episodeRepo.NewEpisodeRepository(dbConnection)

	// AdminPanel Microservice
	grpcConn, err := grpc.Dial(consts.AdminPanelAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()
	adminPanelClient := admin.NewAdminPanelClient(grpcConn)

	// Usecases
	genreUcase := genreUsecase.NewGenreUsecase(genreRepo)
	countryUcase := countryUsecase.NewCountryUsecase(countryRepo, adminPanelClient)
	actorUcase := actorUsecase.NewActorUseCase(actorRepo, adminPanelClient)
	directorUcase := directorUsecase.NewDirectorUseCase(directorRepo, adminPanelClient)
	contentUcase := contentUsecase.NewContentUsecase(contentRepo, countryUcase, genreUcase, actorUcase, directorUcase)
	movieUcase := movieUsecase.NewMovieUsecase(movieRepo, contentUcase)
	tvshowUcase := tvshowUsecase.NewTVShowUsecase(tvshowRepo, contentUcase)
	ratingUcase := ratingUsecase.NewRatingUseCase(ratingRepo, contentUcase)
	favouriteUcase := favouriteUsecase.NewFavouriteUsecase(favouriteRepo)
	seasonUcase := seasonUsecase.NewSeasonUsecase(seasonRepo, tvshowUcase)
	episodeUcase := episodeUsecase.NewEpisodeUsecase(episodeRepo, seasonUcase)
	searchUcase := searchUsecase.NewSearchUsecase(actorRepo, movieRepo, tvshowRepo)

	// Session microservice
	sessionGrpcConn, err := grpc.Dial(consts.SessionblockAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer sessionGrpcConn.Close()
	sessBlockClient := sessionGRPC.NewSessionBlockClient(sessionGrpcConn)
	sessUcase := sessionUsecase.NewSessionUsecase(sessBlockClient)

	// Userblock microservice
	userblockGrpcConn, err := grpc.Dial(consts.UserblockAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer userblockGrpcConn.Close()
	userBlockClient := user.NewUserBlockClient(userblockGrpcConn)
	userUcase := userUsecase.NewUserUsecase(userBlockClient)

	// Monitoring
	e := echo.New()
	mntng := monitoring.NewMonitoring(e)

	// Middleware
	mw := mwares.NewMiddlewareManager(sessUcase, userUcase, mntng)
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
	contentHandler := contentHandler.NewContentHandler(contentUcase)
	movieHandler := movieHandler.NewMovieHandler(movieUcase, contentUcase, countryUcase, genreUcase, actorUcase, directorUcase)
	tvshowHandler := tvshowHandler.NewTVShowHandler(tvshowUcase, contentUcase, countryUcase, genreUcase, actorUcase, directorUcase, seasonUcase)
	ratingHandler := ratingHandler.NewRatingHandler(ratingUcase)
	favouriteHandler := favouriteHandler.NewFavouriteHandler(favouriteUcase, contentUcase)
	seasonHandler := seasonHandler.NewSeasonHandler(seasonUcase)
	episodeHandler := episodeHandler.NewEpisodeHandler(episodeUcase)
	searchHandler := searchHandler.NewSearchHandler(searchUcase)

	userHandler.Configure(e, mw)
	sessionHandler.Configure(e, mw)
	genreHandler.Configure(e, mw)
	countryHandler.Configure(e, mw)
	actorHandler.Configure(e, mw)
	directorHandler.Configure(e, mw)
	contentHandler.Configure(e, mw)
	movieHandler.Configure(e, mw)
	tvshowHandler.Configure(e, mw)
	ratingHandler.Configure(e, mw)
	favouriteHandler.Configure(e, mw)
	seasonHandler.Configure(e, mw)
	episodeHandler.Configure(e, mw)
	searchHandler.Configure(e, mw)

	log.Fatal(e.Start(config.GetServerConnString()))
}
