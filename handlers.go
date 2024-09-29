package main

import (
	"database/sql"
	"log/slog"
	"net/http"
)

type handlers struct {
	db *sql.DB
}

func generateHandlers(db *sql.DB) *handlers {
	return &handlers{db: db}
}

func (h *handlers) handleHome(w http.ResponseWriter, r *http.Request) error {
	reqID := getRequestID(r.Context())
	slog.Info("loggin request", "requestId", reqID)

	writeJSON(w, 200, map[string]string{
		"elo": "yep",
	})
	return nil
}
