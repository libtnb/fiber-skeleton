package bootstrap

import (
	"github.com/libtnb/utils/crypt"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

// NewCrypter builds the crypter for encrypting values at rest.
func NewCrypter(config *conf.Config) (crypt.Crypter, error) {
	return crypt.NewXChacha20Poly1305([]byte(config.App.Key))
}
