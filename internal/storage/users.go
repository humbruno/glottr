package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/humbruno/glottr/internal/env"
	"github.com/humbruno/glottr/internal/queries"
)

const (
	realm             = "glottr"
	temporaryPassword = false
)

var (
	ErrUsersFailedIdpConnection      = errors.New("Failed to connect to IDP")
	ErrUsersFailedCreateUserIdp      = errors.New("Failed to create user in IDP")
	ErrUsersUserExists               = errors.New("User already exists")
	ErrUsersFailedSetIdpPassword     = errors.New("Failed to ser user password in IDP")
	ErrUsersFailedInsertUserDatabase = errors.New("Failed to insert IDP user into database")
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStorage struct {
	idp *gocloak.GoCloak
	db  *sql.DB
}

type Users interface {
	Create(ctx context.Context, usr *User) error
}

func (s *UserStorage) handleIdpAdminLogin(ctx context.Context) (*gocloak.JWT, error) {
	usr := env.GetString("KC_CLI_ADMIN_USERNAME", "")
	pswd := env.GetString("KC_CLI_ADMIN_PASSWORD", "")

	token, err := s.idp.LoginAdmin(ctx, usr, pswd, realm)
	if err != nil {
		slog.Error(ErrUsersFailedIdpConnection.Error(), "err", err)
		return nil, ErrUsersFailedIdpConnection
	}

	return token, nil
}

func (s *UserStorage) createIdpUser(ctx context.Context, usr *User) (userId string, err error) {
	admin, err := s.handleIdpAdminLogin(ctx)
	if err != nil {
		return "", err
	}

	newUser := gocloak.User{
		Email:    gocloak.StringP(usr.Email),
		Username: gocloak.StringP(usr.Username),
		Enabled:  gocloak.BoolP(true),
	}

	newUserId, err := s.idp.CreateUser(ctx, admin.AccessToken, realm, newUser)
	if err != nil {

		if strings.Contains(err.Error(), "409") {
			slog.Error(ErrUsersUserExists.Error(), "err", err)
			return "", ErrUsersUserExists
		}

		slog.Error(ErrUsersFailedCreateUserIdp.Error(), "err", err)
		return "", ErrUsersFailedCreateUserIdp
	}

	slog.Info("User created in IDP", "id", newUserId)

	err = s.setIdpUserPassword(ctx, newUserId, usr.Password)
	if err != nil {
		return "", err
	}

	return newUserId, nil
}

func (s *UserStorage) setIdpUserPassword(ctx context.Context, userId string, pswd string) error {
	admin, err := s.handleIdpAdminLogin(ctx)
	if err != nil {
		return err
	}

	err = s.idp.SetPassword(ctx, admin.AccessToken, userId, realm, pswd, temporaryPassword)
	if err != nil {
		slog.Error(ErrUsersFailedSetIdpPassword.Error(), "err", err)
		return ErrUsersFailedSetIdpPassword
	}

	slog.Info("User password set in IDP")
	return nil
}

func (s *UserStorage) Create(ctx context.Context, usr *User) error {
	id, err := s.createIdpUser(ctx, usr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.createDbUser(ctx, tx, usr, id); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStorage) createDbUser(ctx context.Context, tx *sql.Tx, usr *User, id string) error {
	var createdID string

	err := tx.QueryRowContext(ctx, queries.InsertUser, id, usr.Username, usr.Email).Scan(&createdID)
	if err != nil {
		slog.Error(ErrUsersFailedInsertUserDatabase.Error(), "err", err)
		return ErrUsersFailedInsertUserDatabase
	}

	slog.Info("IDP user added to database", "id", createdID)

	return nil
}
