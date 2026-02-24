package entity

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Item struct {
	ChrtID      int       `json:"chrt_id,omitempty"`
	TrackNumber string    `json:"track_number,omitempty"`
	Price       int       `json:"price,omitempty"`
	RID         uuid.UUID `json:"rid,omitempty"`
	Name        string    `json:"name,omitempty"`
	Sale        int       `json:"sale,omitempty"`
	Size        string    `json:"size,omitempty"`
	TotalPrice  int       `json:"total_price,omitempty"`
	NmID        int       `json:"nm_id,omitempty"`
	Brand       string    `json:"brand,omitempty"`
	Status      int       `json:"status,omitempty"`
}

func (i Item) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ChrtID, validation.Required),
		validation.Field(&i.TrackNumber, validation.Required),
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.RID, validation.Required),
		validation.Field(&i.TotalPrice, validation.Required, validation.Min(1)),
		validation.Field(&i.NmID, validation.Required),
	)
}
