package storage

import (
	"database/sql"
)

type Storage struct {
	Users
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStorage{db},
	}
}
