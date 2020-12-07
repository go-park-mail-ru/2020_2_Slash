package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2020_2_Slash/config"
	grpcSess "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery/grpc"
	sessionRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/session/repository"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
)

func main() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	logger.InitLogger(config.GetLoggerDir(), config.GetLogLevel())

	// Database
	dbConnection, err := sql.Open("postgres", config.GetProdDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	authMsAdress := config.GetAuthMSConnString()
	lis, err := net.Listen("tcp", authMsAdress)
	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	sessRepo := sessionRepo.NewSessionPgRepository(dbConnection)

	server := grpc.NewServer()
	grpcSess.RegisterSessionBlockServer(server, grpcSess.NewSessionBlockMicroservice(sessRepo))

	logger.Println("Starting server at", authMsAdress)
	server.Serve(lis)
}
