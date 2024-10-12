package main

import (
	"net/http"

	"github.com/humbruno/glottr/internal/auth"
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

	hash, err := auth.MakePasswordHash(payload.Password)
	if err != nil {
		app.internalError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"hash": hash,
	})
}
