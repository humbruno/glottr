package database

import (
	"context"
	"database/sql"
	"time"
)

func New(driver, addr string) (*sql.DB, error) {
	db, err := sql.Open(driver, addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
