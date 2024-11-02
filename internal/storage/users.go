package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
	"github.com/humbruno/glottr/internal/env"
)

const (
	realm                    = "glottr"
	fallbackIdpAdminUsername = "bruno"
	fallbackIdpAdminPassword = "admin"
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
	usr := env.GetString("KC_CLI_ADMIN_USERNAME", fallbackIdpAdminUsername)
	pswd := env.GetString("KC_CLI_ADMIN_PASSWORD", fallbackIdpAdminPassword)
	return s.idp.LoginAdmin(ctx, usr, pswd, realm)
}

func (s *UserStorage) CreateUser(ctx context.Context, tx *sql.Tx, email, username string) error {
	admin, err := s.handleIdpAdminLogin(ctx)
	if err != nil {
		slog.Error("Failed to connect as keycloak admin", "err", err)
		return err
	}

	newUser := gocloak.User{
		Email:    gocloak.StringP(email),
		Username: gocloak.StringP(username),
		Enabled:  gocloak.BoolP(true),
	}

	id, err := s.idp.CreateUser(ctx, admin.AccessToken, realm, newUser)
	if err != nil {
		slog.Error("Failed to create user in keycloak", "err", err)
		return err
	}

	slog.Info("created user!", "id", id)

	query := `
    INSERT INTO users (id, username, email) VALUES 
    ($1, $2, $3)
    RETURNING id, created_at
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	tx.QueryRowContext(
		ctx,
		query,
		username,
		email,
	)

	return nil
}
