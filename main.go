package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type app struct {
	db       *sql.DB
	handlers *handlers
	port     int
}

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading environment variables")
		os.Exit(1)
	}

	portEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		slog.Warn("Environment variable 'PORT' cannot be converted to int, defaulting to port 8000", "portEnv", portEnv)
		port = 8000
	}

	dbUrl := os.Getenv("DB_URL")
	dbAuthToken := os.Getenv("DB_AUTH_TOKEN")
	connString := fmt.Sprintf("%s?authToken=%s", dbUrl, dbAuthToken)
	db, err := sql.Open("libsql", connString)
	if err != nil {
		slog.Error("Failed to open db connection", "connectionString", connString, "error", err)
		panic("Failed to open db connection")
	}
	defer db.Close()

	handlers := generateHandlers(db)

	app := app{
		db:       db,
		port:     port,
		handlers: handlers,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", makeHandlerFunc(app.handlers.handleHome))
	muxWithMiddleware := reqIdMiddleware(mux)

	slog.Info("Starting http server", "port", app.port)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", app.port), muxWithMiddleware); err != nil {
		slog.Error("Failed to start http server", "error", err)
		panic("Failed to start http server")
	}
}
