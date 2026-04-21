package calendar

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Event struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}

func (e Event) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.UserID, validation.Required),
		validation.Field(&e.Date, validation.Required),
		validation.Field(&e.Text, validation.Required),
	)
}

type CreateRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}

func (r CreateRequest) ToEvent() Event {
	return Event{
		ID:     uuid.New(),
		UserID: r.UserID,
		Date:   r.Date,
		Text:   r.Text,
	}
}
