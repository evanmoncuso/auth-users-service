package models

import (
	"errors"

	"auth-users-service/pkg/db"
)

// update userinfo set username=$1 where uid=$2
const preparedInsert = `
	INSERT INTO users ( uuid, first_name, last_name, username, password, email_address )
	VALUES ( $1, $2, $3, $4, $5, $6 )
`

const preparedSelectUserByUUID = `
	SELECT uuid, first_name, last_name, username, email_address from users WHERE uuid = $1
`

// User is all the information we'll store about the user
type User struct {
	// UserID       int64     `jsonapi:"primary,users"`
	UserUUID     string `jsonapi:"primary,users"`
	FirstName    string `jsonapi:"attr,firstName,omitempty"`
	LastName     string `jsonapi:"attr,lastName,omitempty"`
	Username     string `jsonapi:"attr,username"`
	Password     string `jsonapi:"attr,password,omitempty"`
	EmailAddress string `jsonapi:"attr,emailAddress,omitempty"`
}

// Create creates a user and saves it to the database
func (u *User) Create() error {
	if u.Username == "" {
		return errors.New("No Username set. Username is required")
	} else if u.Password == "" {
		return errors.New("No Password set. Password is required")
	}

	statement, err := datastore.DB.Prepare(preparedInsert)
	if err != nil {
		return err
	}

	_, err = statement.Exec(
		u.UserUUID,
		u.FirstName,
		u.LastName,
		u.Username,
		u.Password,
		u.EmailAddress,
	)

	if err != nil {
		return err
	}

	return nil
}

// FindRecordByUUID retrieves a user object from the DB
func FindRecordByUUID(uuid string) (User, error) {
	var user User

	err := datastore.DB.QueryRow(preparedSelectUserByUUID, uuid).Scan(
		&user.UserUUID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.EmailAddress,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}
