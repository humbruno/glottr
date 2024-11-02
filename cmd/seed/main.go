package main

import (
	"log"
	"log/slog"

	"github.com/humbruno/glottr/internal/database"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/storage"
	_ "github.com/lib/pq"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:password@localhost:5432/glottr?sslmode=disable")
	db, err := database.New("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Database connection established")

	storage := storage.NewStorage(db)

	slog.Info("Running Seed fn")
	database.Seed(storage, db)
}
