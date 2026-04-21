package calendar

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=store.go -destination=store_mock.go -typed -package=calendar
type EventStore interface {
	Create(ctx context.Context, r CreateRequest) (Event, error)
	Update(ctx context.Context, event Event) (Event, error)
	Delete(ctx context.Context, r DeleteRequest) (Event, error)
	Query(ctx context.Context, query EventQuery) ([]Event, error)
}

type InMemoryEventStore struct {
	events map[uuid.UUID]map[uuid.UUID]Event
	mu     sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[uuid.UUID]map[uuid.UUID]Event),
	}
}

func (s *InMemoryEventStore) Create(ctx context.Context, r CreateRequest) (Event, error) {
	defer s.mu.Unlock()
	s.mu.Lock()

	event := r.ToEvent()

	_, ok := s.events[event.UserID]
	if !ok {
		s.events[event.UserID] = make(map[uuid.UUID]Event)
	}

	s.events[event.UserID][event.ID] = event

	return event, nil
}

func (s *InMemoryEventStore) Update(ctx context.Context, event Event) (Event, error) {
	defer s.mu.Unlock()
	s.mu.Lock()

	userEvents, ok := s.events[event.UserID]
	if !ok {
		return Event{}, ErrUserNotFound
	}

	_, ok = userEvents[event.ID]
	if !ok {
		return Event{}, ErrEventNotFound
	}

	userEvents[event.ID] = event

	return event, nil
}

func (s *InMemoryEventStore) Delete(ctx context.Context, r DeleteRequest) (Event, error) {
	defer s.mu.Unlock()
	s.mu.Lock()

	userEvents, ok := s.events[r.UserID]
	if !ok {
		return Event{}, ErrUserNotFound
	}

	event, ok := userEvents[r.EventID]
	if !ok {
		return Event{}, ErrEventNotFound
	}

	delete(userEvents, r.EventID)

	return event, nil
}

func (s *InMemoryEventStore) Query(ctx context.Context, query EventQuery) ([]Event, error) {
	defer s.mu.RUnlock()
	s.mu.RLock()

	userEvents, ok := s.events[query.UserID]
	if !ok {
		return nil, ErrUserNotFound
	}

	now := time.Now().UTC()

	var predicate func(Event) bool

	switch query.Period {
	case PeriodNone:
		events := slices.Collect(maps.Values(userEvents))

		slices.SortFunc(events, func(a, b Event) int {
			return a.Date.Compare(b.Date)
		})

		return events, nil
	case PeriodDay:
		predicate = func(e Event) bool {
			date := e.Date.UTC()

			return now.Year() == date.Year() &&
				now.YearDay() == date.YearDay()
		}
	case PeriodWeek:
		predicate = func(e Event) bool {
			year, week := now.ISOWeek()
			eventYear, eventWeek := e.Date.ISOWeek()

			return year == eventYear &&
				week == eventWeek
		}
	case PeriodMonth:
		predicate = func(e Event) bool {
			date := e.Date.UTC()

			return now.Year() == date.Year() &&
				now.Month() == date.Month()
		}
	default:
		return nil, ErrInvalidPeriod
	}

	var filtered []Event

	for _, event := range userEvents {
		if predicate(event) {
			filtered = append(filtered, event)
		}
	}

	slices.SortFunc(filtered, func(a, b Event) int {
		return a.Date.Compare(b.Date)
	})

	return filtered, nil
}
