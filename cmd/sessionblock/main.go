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
	dbConnection, err := sql.Open("postgres", config.GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalln("Can't listen port :8082", err)
	}
	defer lis.Close()

	sessRepo := sessionRepo.NewSessionPgRepository(dbConnection)

	server := grpc.NewServer()
	grpcSess.RegisterSessionBlockServer(server, grpcSess.NewSessionBlockMicroservice(sessRepo))

	logger.Println("Starting server at :8082")
	server.Serve(lis)
}
