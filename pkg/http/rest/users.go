package rest

import (
	"auth-users-service/pkg/middleware"
	"auth-users-service/pkg/models"

	"net/http"

	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
)

// CreateUserHandler is responsible for getting the body, saving the request to the database, and returning the created user
func CreateUserHandler(res http.ResponseWriter, req *http.Request) {
	middleware.EnableCors(&res)

	user := new(models.User)

	if err := jsonapi.UnmarshalPayload(req.Body, user); err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// username and password are required
	if user.Username == "" {
		handleHTTPError(res, "No username on request. Username is required", http.StatusBadRequest)
		return
	} else if user.Password == "" {
		handleHTTPError(res, "No password on request. Password is required", http.StatusBadRequest)
		return
	}

	// jsonapi package doesn't like ids as type uuid.UUID so force to string
	user.UserUUID = uuid.New().String()

	pwBytes := []byte(user.Password)
	cryptPw, err := bcrypt.GenerateFromPassword(pwBytes, 4)

	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = string(cryptPw)

	err = user.Create()

	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(res, user); err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetUserHandler returns information on an individual user if the request is authenticated
// authentication set in server/main.go
func GetUserHandler(res http.ResponseWriter, req *http.Request) {
	middleware.EnableCors(&res)
	vars := mux.Vars(req)
	userUUID := vars["uuid"]

	user, err := models.FindRecordByUUID(userUUID)
	if err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(res, &user); err != nil {
		handleHTTPError(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
