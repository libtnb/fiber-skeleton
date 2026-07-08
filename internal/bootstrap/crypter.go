package bootstrap

import (
	"github.com/libtnb/utils/crypt"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/config"
)

func NewCrypter(i do.Injector) (crypt.Crypter, error) {
	return crypt.NewXChacha20Poly1305([]byte(do.MustInvoke[*config.Config](i).App.Key))
}
