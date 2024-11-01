package storage

import (
	"context"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
)

const realm = "glottr"

type UserInfo struct {
	Name  string
	Email string
}

type UserStorage struct {
	idp *gocloak.GoCloak
}

type Users interface {
	CreateUser(ctx context.Context, email, username string) error
}

func (s *UserStorage) CreateUser(ctx context.Context, email, username string) error {
	token, err := s.idp.LoginAdmin(ctx, "bruno", "admin", realm)
	if err != nil {
		slog.Error("Failed to connect as keycloak admin", "err", err)
		return err
	}

	newUser := gocloak.User{
		Email:    gocloak.StringP(email),
		Username: gocloak.StringP(username),
		Enabled:  gocloak.BoolP(true),
	}

	id, err := s.idp.CreateUser(ctx, token.AccessToken, realm, newUser)
	if err != nil {
		slog.Error("Failed to create user in keycloak", "err", err)
		return err
	}

	slog.Info("created user!", "id", id)

	return nil
}
