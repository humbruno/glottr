package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/queries"
)

const (
	realm = "glottr"
)

var (
	ErrUsersFailedIdpConnection      = errors.New("Failed to connect to IDP")
	ErrUsersFailedCreateUserIdp      = errors.New("Failed to create user in IDP")
	ErrUsersFailedInsertUserDatabase = errors.New("Failed to insert IDP user into database")
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

func (s *UserStorage) CreateUser(ctx context.Context, tx *sql.Tx, usr User) error {
	admin, err := s.handleIdpAdminLogin(ctx)
	if err != nil {
		slog.Error(ErrUsersFailedIdpConnection.Error(), "err", err)
		return ErrUsersFailedIdpConnection
	}

	newUser := gocloak.User{
		Email:    gocloak.StringP(usr.Email),
		Username: gocloak.StringP(usr.Username),
		Enabled:  gocloak.BoolP(true),
	}

	id, err := s.idp.CreateUser(ctx, admin.AccessToken, realm, newUser)
	if err != nil {
		slog.Error(ErrUsersFailedCreateUserIdp.Error(), "err", err)
		return ErrUsersFailedCreateUserIdp
	}

	slog.Info("User created in IDP", "id", id)

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var createdID string

	err = tx.QueryRowContext(ctx, queries.InsertUser, id, usr.Username, usr.Email).Scan(&createdID)
	if err != nil {
		slog.Error(ErrUsersFailedInsertUserDatabase.Error(), "err", err)
		return ErrUsersFailedInsertUserDatabase
	}

	slog.Info("IDP user added to database", "id", createdID)

	return nil
}
