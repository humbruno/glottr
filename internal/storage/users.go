package storage

import (
	"context"
	"database/sql"
)

type UserStorage struct {
	db *sql.DB
}

type Users interface {
	GetByID(ctx context.Context, id int) error
}

func (s *UserStorage) GetByID(ctx context.Context, id int) error {
	return nil
}
