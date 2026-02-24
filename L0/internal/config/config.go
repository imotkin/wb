package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/imotkin/L0/internal/api/handler"
	"github.com/imotkin/L0/internal/api/server"
	"github.com/imotkin/L0/internal/broker"
	"github.com/imotkin/L0/internal/cache"
	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/repo/postgres"
)

type Config struct {
	Server   *server.Config   `koanf:"server"`
	Postgres *postgres.Config `koanf:"postgres"`
	Logging  *logger.Config   `koanf:"logging"`
	Broker   *broker.Config   `koanf:"broker"`
	Web      *handler.Config  `koanf:"web"`
	Cache    *cache.Config    `koanf:"cache"`
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
		validation.Field(&c.Postgres, validation.Required),
		validation.Field(&c.Logging, validation.Required),
		validation.Field(&c.Broker, validation.Required),
		validation.Field(&c.Web, validation.Required),
		validation.Field(&c.Cache, validation.Required),
	)
}
