package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_Slash/config"
	actorRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/repository"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminGrpc "github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/content/repository"
	countryRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/country/repository"
	directorRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/director/repository"
	episodeRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/episode/repository"
	genreRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/repository"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	movieRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/repository"
	seasonRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/season/repository"
	tvshowRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/repository"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
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

	lis, err := net.Listen("tcp", consts.AdminPanelAddress)
	if err != nil {
		log.Fatalln("Can't listen address"+consts.AdminPanelAddress, err)
	}
	defer lis.Close()

	genreRepo := genreRepo.NewGenrePgRepository(dbConnection)
	countryRepo := countryRepo.NewCountryPgRepository(dbConnection)
	actorRepo := actorRepo.NewActorPgRepository(dbConnection)
	directorRepo := directorRepo.NewDirectorPgRepository(dbConnection)
	contentRepo := contentRepo.NewContentPgRepository(dbConnection)
	movieRepo := movieRepo.NewMoviePgRepository(dbConnection)
	tvshowRepo := tvshowRepo.NewTVShowPgRepository(dbConnection)
	seasonRepo := seasonRepo.NewSeasonPgRepository(dbConnection)
	episodeRepo := episodeRepo.NewEpisodeRepository(dbConnection)

	server := grpc.NewServer()
	admin.RegisterAdminPanelServer(server, adminGrpc.NewAdminMicroservice(actorRepo,
		directorRepo, countryRepo, genreRepo, movieRepo,
		contentRepo, seasonRepo, episodeRepo, tvshowRepo))

	fmt.Println("Starting server at " + consts.AdminPanelAddress)
	server.Serve(lis)
}
