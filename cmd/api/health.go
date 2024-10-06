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
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalError(w, r, err)
	}
}
