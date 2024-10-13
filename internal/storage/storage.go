package storage

import (
	"database/sql"

	"github.com/Nerzal/gocloak/v13"
)

type Storage struct {
	Users
}

func NewStorage(db *sql.DB, idp *gocloak.GoCloak) Storage {
	return Storage{
		Users: &UserStorage{idp},
	}
}
