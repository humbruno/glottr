package main

import (
	"database/sql"
	"net/http"
)

type handlers struct {
	db *sql.DB
}

func generateHandlers(db *sql.DB) *handlers {
	return &handlers{db: db}
}

func (h *handlers) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Yep"))
}
