package postgres

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	User          string `koanf:"user"`
	Password      string `koanf:"password"`
	Host          string `koanf:"host"`
	Port          string `koanf:"port"`
	Database      string `koanf:"database"`
	ModeSSL       string `koanf:"mode_ssl"`
	MigrationsDir string `koanf:"migrations_dir"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.User, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.Database, validation.Required),
		validation.Field(&c.ModeSSL, validation.Required),
		validation.Field(&c.MigrationsDir, validation.Required),
	)
}

func (c *Config) ConnectionURL() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.User, c.Password, c.ModeSSL,
	)
}
