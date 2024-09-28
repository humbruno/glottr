package main

import (
	"database/sql"
	"fmt"
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	portEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Printf("Port %s cannot be converted to int", portEnv)
		port = 8000
	}

	dbUrl := os.Getenv("DB_URL")
	dbAuthToken := os.Getenv("DB_AUTH_TOKEN")
	connString := fmt.Sprintf("%s?authToken=%s", dbUrl, dbAuthToken)
	db, err := sql.Open("libsql", connString)
	if err != nil {
		log.Fatalf("Failed to open db %s: %s", connString, err)
	}

	app := app{
		db:       db,
		port:     port,
		handlers: &handlers{},
	}
	defer app.db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handlers.handleHome)

	log.Printf("Listening on port %d \n", app.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.port), mux)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
