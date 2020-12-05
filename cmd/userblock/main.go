package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	mygrpc "github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery/grpc"
	userRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/user/repository"
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

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Can't listen port :8081", err)
	}
	defer lis.Close()

	userRepo := userRepo.NewUserPgRepository(dbConnection)

	server := grpc.NewServer()
	mygrpc.RegisterUserBlockServer(server, mygrpc.NewUserblockMicroservice(userRepo))

	fmt.Println("Starting server at :8081")
	server.Serve(lis)
}
