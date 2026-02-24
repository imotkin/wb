package logger

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Format Format `koanf:"format"`
	Level  Level  `koanf:"level"`
	Output Output `koanf:"output"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Format, validation.Required),
		validation.Field(&c.Level, validation.Required),
		validation.Field(&c.Output, validation.Required),
	)
}
