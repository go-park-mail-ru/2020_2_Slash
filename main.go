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

	fmt.Println("Starting server at :8000")
	http.ListenAndServe(":8000", router)
}
