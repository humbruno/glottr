package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) logInternalError(r *http.Request, err error) {
	reqId := middleware.GetReqID(r.Context())
	slog.Error("Internal Error", "requestId", reqId, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "error", err.Error())
}
