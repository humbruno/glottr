package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/humbruno/glottr/internal/database"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/storage"
	"github.com/joho/godotenv"
)

const (
	fallbackLocalDbUrl    = "postgres://admin:adminpassword@localhost/glottr?sslmode=disable"
	fallbackLocalDbDriver = "postgres"
	fallbackListenAddr    = ":8000"
	fallbackLocalhostUrl  = "localhost:8000"
	fallbackEnv           = "DEVELOPMENT"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	cfg := config{
		addr:   env.GetString("ADDR", fallbackListenAddr),
		apiUrl: env.GetString("EXTERNAL_URL", fallbackLocalhostUrl),
		env:    env.GetString("ENV", fallbackEnv),
		db: dbConfig{
			driver: env.GetString("DB_DRIVER", fallbackLocalDbDriver),
			addr:   env.GetString("DB_ADDR", fallbackLocalDbUrl),
		},
	}

	db, err := database.New(cfg.db.driver, cfg.db.addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Database connection established")

	app := application{
		config:  cfg,
		storage: storage.NewStorage(db),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
