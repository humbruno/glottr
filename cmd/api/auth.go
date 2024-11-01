package main

import (
	"net/http"
)

type registerUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.internalError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"hey": "alo",
	})
}
