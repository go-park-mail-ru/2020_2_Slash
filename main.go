package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_Slash/app/handlers"
	"github.com/go-park-mail-ru/2020_2_Slash/app/helpers"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"net/http"
)

const origin = "http://www.flicksbox.ru"

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(echo.HeaderAccessControlAllowOrigin, origin)
		w.Header().Add(echo.HeaderAccessControlAllowMethods, "DELETE, PUT, POST, GET, OPTIONS")
		w.Header().Add(echo.HeaderAccessControlAllowHeaders, "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

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

	siteRouter := EnableCORS(router)

	fmt.Println("Starting server at :8000")
	http.ListenAndServe(":8000", siteRouter)
}
