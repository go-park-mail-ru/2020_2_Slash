package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"

	"github.com/go-park-mail-ru/2020_2_Slash/config"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"

	sessionRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/session/repository"
	sessionUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/session/usecases"

	sessionHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery"
	userHandler "github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery"
	userRepo "github.com/go-park-mail-ru/2020_2_Slash/internal/user/repository"
	userUsecase "github.com/go-park-mail-ru/2020_2_Slash/internal/user/usecases"
)

func main() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Locale storage
	avatarsPath := config.GetAvatarsPath()
	helpers.InitAvatarStorage(avatarsPath)

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

	// Usecases
	sessUcase := sessionUsecase.NewSessionUsecase(sessRepo)
	userUcase := userUsecase.NewUserUsecase(userRepo)

	// Middleware
	e := echo.New()
	mw := mwares.NewMiddlewareManager(sessUcase)
	e.Use(mw.PanicRecovering, mw.AccessLog, mw.CORS)

	e.Static("/avatars", avatarsPath)

	// Delivery
	userHandler := userHandler.NewUserHandler(userUcase, sessUcase)
	sessionHandler := sessionHandler.NewSessionHandler(sessUcase, userUcase)
	userHandler.Configure(e, mw)
	sessionHandler.Configure(e, mw)

	log.Fatal(e.Start(config.GetServerConnString()))
}
