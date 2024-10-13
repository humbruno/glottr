package storage

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type UserInfo struct {
	Name  string
	Email string
}

type UserStorage struct {
	idp *gocloak.GoCloak
}

type Users interface {
	GetInfo(ctx context.Context, accessToken, realm string) (*UserInfo, error)
}

func (s *UserStorage) GetInfo(ctx context.Context, accessToken, realm string) (*UserInfo, error) {
	return nil, nil
}
