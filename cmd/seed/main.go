package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/humbruno/glottr/internal/database"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	addr := env.GetString("DB_ADDR", "postgres://admin:password@localhost:5432/glottr?sslmode=disable")
	db, err := database.New("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Database connection established")

	storage := storage.NewStorage(db)

	slog.Info("Running seed fn")
	database.Seed(storage, db)
}
