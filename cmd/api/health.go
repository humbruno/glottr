package main

import (
	"net/http"
)

// @Summary		Returns the current health status of the application
// @Description	Healthcheck endpoint
// @Tags			Healthcheck
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/v1/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJSON(w, http.StatusOK, nil); err != nil {
		app.logInternalError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
