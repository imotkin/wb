package calendar

import "github.com/google/uuid"

type Period int

const (
	PeriodNone Period = iota
	PeriodDay
	PeriodWeek
	PeriodMonth
)

type EventQuery struct {
	UserID uuid.UUID
	Period Period
}

type DeleteRequest struct {
	EventID uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
}
