package main

import (
	"net/http"

	"github.com/humbruno/glottr/internal/storage"
)

// @Summary		Registers user
// @Description	Register user in the IDP and DB
// @Tags			Register
// @Produce		json
// @Success		200
// @Router			/v1/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userPayload *storage.User

	if err := readJSON(w, r, &userPayload); err != nil {
		app.internalError(w, r, err)
		return
	}

	err := app.storage.Users.Create(r.Context(), userPayload)
	if err != nil {
		app.jsonResponse(w, http.StatusOK, map[string]string{
			"message": "failed to create user",
			"error":   err.Error(),
		})
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"msg": "user created!",
	})
}

// @Summary		Logs in user
// @Description Logs in user if already exists
// @Tags			Login
// @Produce		json
// @Success		200
// @Router			/v1/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	app.jsonResponse(w, http.StatusOK, map[string]string{
		"msg": "logged in",
	})
}
