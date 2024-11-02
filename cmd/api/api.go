package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/humbruno/glottr/docs"
	"github.com/humbruno/glottr/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	storage storage.Storage
	config  config
}

type config struct {
	addr   string
	apiUrl string
	env    string
	db     dbConfig
}

type dbConfig struct {
	driver string
	addr   string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Post("/register", app.registerUserHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Host = app.config.apiUrl

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("Starting http server", "port", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	slog.Info("Server has stopped", "addr", app.config.addr)

	return nil
}
