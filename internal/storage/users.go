package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
	"github.com/humbruno/glottr/internal/env"
)

const (
	realm = "glottr"
)

type User struct {
	Username string
	Email    string
}

type UserStorage struct {
	idp *gocloak.GoCloak
}

type Users interface {
	CreateUser(ctx context.Context, tx *sql.Tx, email, username string) error
}

func (s *UserStorage) handleIdpAdminLogin(ctx context.Context) (*gocloak.JWT, error) {
	usr := env.GetString("KC_CLI_ADMIN_USERNAME", "")
	pswd := env.GetString("KC_CLI_ADMIN_PASSWORD", "")
	return s.idp.LoginAdmin(ctx, usr, pswd, realm)
}

func (s *UserStorage) CreateUser(ctx context.Context, tx *sql.Tx, email, username string) error {
	admin, err := s.handleIdpAdminLogin(ctx)
	if err != nil {
		slog.Error("Failed to connect to IDP as admin", "err", err)
		return err
	}

	newUser := gocloak.User{
		Email:    gocloak.StringP(email),
		Username: gocloak.StringP(username),
		Enabled:  gocloak.BoolP(true),
	}

	id, err := s.idp.CreateUser(ctx, admin.AccessToken, realm, newUser)
	if err != nil {
		slog.Error("Failed to create user in IDP", "err", err)
		return err
	}

	slog.Info("User created in IDP", "id", id)

	query := `
    INSERT INTO users (id, username, email) VALUES 
    ($1, $2, $3)
    RETURNING id
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var createdID string

	err = tx.QueryRowContext(ctx, query, id, username, email).Scan(&createdID)
	if err != nil {
		slog.Error("Failed to insert user into database", "err", err)
		return err
	}

	slog.Info("IDP user added to database", "id", createdID)

	return nil
}
