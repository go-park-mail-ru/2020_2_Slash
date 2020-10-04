package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
)

func main() {
	router := mux.NewRouter()
	UserHandler := handlers.NewUserHandler()

	router.HandleFunc("/api/v1/user/register", UserHandler.Register).Methods("POST")
	router.HandleFunc("/api/v1/user/profile", UserHandler.GetUserProfile).Methods("GET")
	router.HandleFunc("/api/v1/user/profile", UserHandler.ChangeUserProfile).Methods("PUT")
	router.HandleFunc("/api/v1/user/login", UserHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/user/logout", UserHandler.Logout).Methods("DELETE")

	fmt.Println("Starting server at :8000")
	http.ListenAndServe(":8000", router)
}
