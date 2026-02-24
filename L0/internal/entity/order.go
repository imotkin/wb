package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Order struct {
	UID               uuid.UUID `json:"order_uid,omitempty"`
	TrackNumber       string    `json:"track_number,omitempty"`
	Entry             string    `json:"entry,omitempty"`
	Delivery          Delivery  `json:"delivery,omitzero"`
	Payment           Payment   `json:"payment,omitzero"`
	Items             []Item    `json:"items,omitempty"`
	Locale            string    `json:"locale,omitempty"`
	InternalSignature string    `json:"internal_signature,omitempty"`
	CustomerID        string    `json:"customer_id,omitempty"`
	DeliveryService   string    `json:"delivery_service,omitempty"`
	ShardKey          string    `json:"shardkey,omitempty"`
	SmID              int       `json:"sm_id,omitempty"`
	DateCreated       time.Time `json:"date_created,omitzero"`
	Shard             string    `json:"oof_shard,omitempty"`
}

func (o Order) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.UID, validation.Required),
		validation.Field(&o.TrackNumber, validation.Required),
		validation.Field(&o.Entry, validation.Required),
		validation.Field(&o.Delivery, validation.Required),
		validation.Field(&o.Payment, validation.Required),
		validation.Field(&o.Items, validation.Required),
		validation.Field(&o.Locale, validation.Required),
		validation.Field(&o.InternalSignature, validation.Required),
		validation.Field(&o.CustomerID, validation.Required),
		validation.Field(&o.DeliveryService, validation.Required),
		validation.Field(&o.ShardKey, validation.Required),
		validation.Field(&o.SmID, validation.Required),
		validation.Field(&o.DateCreated, validation.Required),
		validation.Field(&o.Shard, validation.Required),
	)
}
