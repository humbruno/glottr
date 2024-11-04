package main

import (
	"net/http"

	"github.com/humbruno/glottr/internal/storage"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload *storage.User

	if err := readJSON(w, r, &userPayload); err != nil {
		app.internalError(w, r, err)
		return
	}

	err := app.storage.Users.Create(r.Context(), userPayload)
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
