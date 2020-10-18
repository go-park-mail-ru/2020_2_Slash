package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/helpers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const origin = "http://www.flicksbox.ru"

const (
	host     = "localhost"
	port     = 5432
	user     = "slash"
	password = "slashpassword"
	dbname   = "slash"
)

func getDatabaseConnectionString() string {
	return fmt.Sprintf("host=%s port=%d "+
		"user=%s password=%s dbname=%s sslmode=disable", host,
		port, user, password, dbname)
}

func main() {
	helpers.InitAvatarStorage()

	db, err := sql.Open("postgres", getDatabaseConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database!")
	UserHandler := handlers.NewUserHandler(db)

	router := echo.New()

	router.Use(middleware.BodyLimit("10M"))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}))

	router.Static("/avatars/", "./avatars")

	router.POST("/api/v1/user/register", UserHandler.Register)
	router.GET("/api/v1/user/profile", UserHandler.GetUserProfile)
	router.PUT("/api/v1/user/profile", UserHandler.ChangeUserProfile)
	router.POST("/api/v1/user/login", UserHandler.Login)
	router.DELETE("/api/v1/user/logout", UserHandler.Logout)
	router.GET("/api/v1/session", UserHandler.CheckSession)
	router.POST("/api/v1/user/avatar", UserHandler.SetAvatar)

	router.Logger.Fatal(router.Start(":8000"))
}
