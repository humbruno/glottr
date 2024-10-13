package auth

import (
	"github.com/Nerzal/gocloak/v13"
)

func NewIdpClient(basePath string) *gocloak.GoCloak {
	return gocloak.NewClient(basePath)
}
