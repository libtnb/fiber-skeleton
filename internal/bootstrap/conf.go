package bootstrap

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func NewConf() (*koanf.Koanf, error) {
	conf := koanf.New(".")
	if err := conf.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
		return nil, err
	}

	return conf, nil
}
