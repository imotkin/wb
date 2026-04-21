package http

import "github.com/google/uuid"

type DeleteEventRequest struct {
	ID uuid.UUID `json:"id"`
}
