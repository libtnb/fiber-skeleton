package bootstrap

import (
	"github.com/knadh/koanf/v2"
	"github.com/libtnb/utils/crypt"
)

func NewCrypter(conf *koanf.Koanf) (crypt.Crypter, error) {
	return crypt.NewXChacha20Poly1305([]byte(conf.MustString("app.key")))
}
