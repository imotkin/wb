package calendar

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -typed -package=calendar
type Service interface {
	CreateEvent(ctx context.Context, r CreateRequest) (Event, error)
	UpdateEvent(ctx context.Context, event Event) (Event, error)
	DeleteEvent(ctx context.Context, r DeleteRequest) (Event, error)
	DayEvents(ctx context.Context, userID uuid.UUID) ([]Event, error)
	WeekEvents(ctx context.Context, userID uuid.UUID) ([]Event, error)
	MonthEvents(ctx context.Context, userID uuid.UUID) ([]Event, error)
	Events(ctx context.Context, userID uuid.UUID) ([]Event, error)
}

type EventService struct {
	events EventStore
}

func NewService(store EventStore) *EventService {
	return &EventService{events: store}
}

func (s *EventService) CreateEvent(ctx context.Context, r CreateRequest) (Event, error) {
	return s.events.Create(ctx, r)
}

func (s *EventService) UpdateEvent(ctx context.Context, event Event) (Event, error) {
	return s.events.Update(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, r DeleteRequest) (Event, error) {
	return s.events.Delete(ctx, r)
}

func (s *EventService) DayEvents(ctx context.Context, userID uuid.UUID) ([]Event, error) {
	return s.events.Query(ctx, EventQuery{
		UserID: userID,
		Period: PeriodDay,
	})
}

func (s *EventService) WeekEvents(ctx context.Context, userID uuid.UUID) ([]Event, error) {
	return s.events.Query(ctx, EventQuery{
		UserID: userID,
		Period: PeriodWeek,
	})
}

func (s *EventService) MonthEvents(ctx context.Context, userID uuid.UUID) ([]Event, error) {
	return s.events.Query(ctx, EventQuery{
		UserID: userID,
		Period: PeriodMonth,
	})
}

func (s *EventService) Events(ctx context.Context, userID uuid.UUID) ([]Event, error) {
	return s.events.Query(ctx, EventQuery{
		UserID: userID,
		Period: PeriodNone,
	})
}
