package main

import (
	"net/http"
)

type registerUserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.internalError(w, r, err)
		return
	}

	err := app.storage.Users.CreateUser(r.Context(), payload.Email, payload.Username)
	if err != nil {
		app.jsonResponse(w, http.StatusOK, map[string]string{
			"msg": "failed to create user",
			"err": err.Error(),
		})
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"msg": "user created!",
	})
}
