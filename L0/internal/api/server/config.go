package server

import (
	"net"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Host         string        `koanf:"host"`
	Port         string        `koanf:"port"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
	IdleTimeout  time.Duration `koanf:"idle_timeout"`
}

func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.ReadTimeout, validation.Required),
		validation.Field(&c.WriteTimeout, validation.Required),
		validation.Field(&c.IdleTimeout, validation.Required),
	)
}
