package datastore

import (
	"database/sql"
	"fmt"

	// sweet psql driver

	_ "github.com/lib/pq"
)

// DB is the sharable pool of connections for the app
var DB *sql.DB

// InitializeDB starts a connection to the database and exposes it to the application
func InitializeDB(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	DB = db

	fmt.Println("Successful db connect")
	return nil
}
