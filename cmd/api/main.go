package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/humbruno/glottr/internal/auth"
	"github.com/humbruno/glottr/internal/database"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	fallbackLocalDbName      = "glottr"
	fallbackLocalDbUser      = "admin"
	fallbackLocalDbPassword  = "password"
	fallbackListenAddr       = ":8000"
	fallbackLocalhostUrl     = "localhost:8000"
	fallbackEnv              = "DEVELOPMENT"
	fallbackIdpBaseUrl       = "http://localhost:8080"
	fallbackIdpAdminUsername = "brunoadmin"
	fallbackIdpAdminPassword = "admin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	dbName := env.GetString("POSTGRES_DB", fallbackLocalDbName)
	dbUser := env.GetString("POSTGRES_USER", fallbackLocalDbUser)
	dbPassword := env.GetString("POSTGRES_PASSWORD", fallbackLocalDbPassword)

	cfg := config{
		addr:   env.GetString("ADDR", fallbackListenAddr),
		apiUrl: env.GetString("EXTERNAL_URL", fallbackLocalhostUrl),
		env:    env.GetString("ENV", fallbackEnv),
		db: dbConfig{
			driver: "postgres",
			addr:   fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPassword, dbName),
		},
		idp: idpConfig{
			baseUrl:       env.GetString("KC_BASE_URL", fallbackIdpBaseUrl),
			adminUsername: env.GetString("KC_BOOTSTRAP_ADMIN_USERNAME", fallbackIdpAdminUsername),
			adminPassword: env.GetString("KC_BOOTSTRAP_ADMIN_PASSWORD", fallbackIdpAdminPassword),
		},
	}

	db, err := database.New(cfg.db.driver, cfg.db.addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	slog.Info("Database connection established")

	idp := auth.NewIdpClient(cfg.idp.baseUrl)

	app := application{
		config:  cfg,
		storage: storage.NewStorage(db, idp),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
