package main

import (
	"auth-users-service/pkg/db"
	"auth-users-service/pkg/http/rest"
	"auth-users-service/pkg/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	dbConnectionURL := os.Getenv("DATABASE_URL")

	if port == "" {
		panic("No PORT specified for app")
	}

	err := datastore.InitializeDB(dbConnectionURL)
	if err != nil {
		panic(err)
	}

	connection := fmt.Sprintf("127.0.0.1:%s", port)

	router := mux.NewRouter()
	router.Use(middleware.Logger)

	router.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()

	// create user (NOT AUTHENTICATED)
	api.HandleFunc("/users", rest.CreateUserHandler).Methods("POST")

	// router.Use(middleware.AuthenticateJWT)
	api.HandleFunc("/users/{uuid}", rest.GetUserHandler).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         connection,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
