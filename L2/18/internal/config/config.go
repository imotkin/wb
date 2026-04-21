package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imotkin/L2/18/internal/http"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server *http.ServerConfig `koanf:"server"`
}

func Parse(path string) (*Config, error) {
	var (
		k   = koanf.New(".")
		cfg = new(Config)
	)

	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}

	err = k.Unmarshal("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Server, validation.Required),
	)
}
