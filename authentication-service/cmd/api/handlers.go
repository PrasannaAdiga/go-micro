package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r http.Request) {

	var requestPaylod struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPaylod)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate use against the DB
	user, err := app.Models.User.GetByEmail(requestPaylod.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials!"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPaylod.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid credentials!"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
