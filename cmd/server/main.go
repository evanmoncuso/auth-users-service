package main

import (
	"auth-validation/pkg/validate"

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
	tokenConnectionURL := os.Getenv("REDIS_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	if port == "" {
		panic("No PORT specified for app")
	}

	err := datastore.InitializeDB(dbConnectionURL)
	if err != nil {
		panic(err)
	}

	err = validate.Initialize(tokenConnectionURL, tokenSecret)
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

	protected := api.PathPrefix("").Subrouter()
	protected.Use(validate.AuthenticateRoute)

	// create user (NOT AUTHENTICATED)
	api.HandleFunc("/users", rest.CreateUserHandler).Methods("POST")

	// get user info AUTHENTICATED
	protected.HandleFunc("/users/{uuid}", rest.GetUserHandler).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         connection,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on port: %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
