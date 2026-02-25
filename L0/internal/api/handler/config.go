package handler

import (
	"errors"
	"fmt"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	TemplatePath string `koanf:"template_path"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(
			&c.TemplatePath,
			validation.Required.Error("template_path is required"),
			validation.By(func(value any) error {
				path, ok := value.(string)
				if !ok {
					return errors.New("template_path is not a string")
				}

				info, err := os.Stat(path)
				if err != nil {
					return fmt.Errorf("get file info: %w", err)
				}

				parts := strings.Split(info.Name(), ".")
				if len(parts) != 2 {
					return fmt.Errorf("invalid template file: %s", info.Name())
				}

				if ext := parts[1]; ext != "html" {
					return fmt.Errorf("invalid template extension: %s", ext)
				}

				return nil
			}),
		),
	)
}
