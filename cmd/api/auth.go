package main

import "net/http"

type registerUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.logError(r, err)
		return
	}
}
