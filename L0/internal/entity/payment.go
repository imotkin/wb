package entity

import "github.com/google/uuid"

type Payment struct {
	Transaction  uuid.UUID `json:"transaction,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	Currency     string    `json:"currency,omitempty"`
	Provider     string    `json:"provider,omitempty"`
	Amount       int       `json:"amount,omitempty"`
	PaymentDt    int       `json:"payment_dt,omitempty"`
	Bank         string    `json:"bank,omitempty"`
	DeliveryCost int       `json:"delivery_cost,omitempty"`
	GoodsTotal   int       `json:"goods_total,omitempty"`
	CustomFee    int       `json:"custom_fee,omitempty"`
}
