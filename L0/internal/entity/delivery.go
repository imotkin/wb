package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Delivery struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Zip     string `json:"zip,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
	Region  string `json:"region,omitempty"`
	Email   string `json:"email,omitempty"`
}

func (d Delivery) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(
			&d.Name,
			validation.Required,
		),
		validation.Field(
			&d.Phone,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\+\d{11}$`)),
		),
		validation.Field(
			&d.Zip,
			validation.Required,
			validation.Match(regexp.MustCompile("^[0-9]{6}$")),
		),
		validation.Field(
			&d.City,
			validation.Required,
		),
		validation.Field(
			&d.Address,
			validation.Required,
		),
		validation.Field(
			&d.Region,
			validation.Required,
		),
		validation.Field(
			&d.Email,
			validation.Required,
			is.Email,
		),
	)
}
