package storage

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/humbruno/glottr/internal/auth"
	"github.com/humbruno/glottr/internal/env"
)

const fallbackIdpBaseUrl = "http://localhost:8080"

var queryTimeoutDuration = time.Second * 5

type Storage struct {
	Users
}

func NewStorage(db *sql.DB) Storage {
	idpUrl := env.GetString("KC_BASE_URL", fallbackIdpBaseUrl)
	idp := auth.NewIdpClient(idpUrl)
	slog.Info("Idp connection established")

	return Storage{
		Users: &UserStorage{idp},
	}
}
