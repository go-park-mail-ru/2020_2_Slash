package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/helpers"
)

func main() {
	helpers.InitAvatarStorage()
	router := mux.NewRouter()
	UserHandler := handlers.NewUserHandler()

	router.PathPrefix("/avatars/").Handler(http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars")))) // TODO: config
	router.HandleFunc("/api/v1/session", UserHandler.CheckSession).Methods("GET")
	router.HandleFunc("/api/v1/user/register", UserHandler.Register).Methods("POST")
	router.HandleFunc("/api/v1/user/profile", UserHandler.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/user/profile", UserHandler.ChangeUserProfile).Methods("PUT")
	router.HandleFunc("/api/v1/user/login", UserHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/user/logout", UserHandler.Logout).Methods("DELETE")
	router.HandleFunc("/api/v1/user/avatar", UserHandler.SetAvatar).Methods("POST")

	fmt.Println("Starting server at :8000")
	http.ListenAndServe(":8000", router)
}
