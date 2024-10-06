package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) logError(r *http.Request, err error) {
	reqId := middleware.GetReqID(r.Context())
	slog.Error("Err", "requestId", reqId, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "error", err.Error())
}

func (app *application) internalError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	writeJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) unprocessableEntity(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	writeJSONError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
}
