package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matryer/is"
)

func weekRange(t time.Time) []time.Time {
	weekday := int(t.Weekday() % 7)

	dates := make([]time.Time, 0, 7)

	for i := -weekday + 1; i <= (7 - weekday); i++ {
		offset := time.Duration(i) * time.Hour * 24
		dates = append(dates, t.Add(offset).Round(0))
	}

	return dates
}

func monthRange(t time.Time) []time.Time {
	var (
		year  = t.Year()
		month = t.Month()
		dates = make([]time.Time, 0, 31)
	)

	for day := range 31 {
		date := time.Date(
			year, month, day+1,
			10, 00, 00, 00, time.Local,
		)

		if date.Month() != month {
			break
		}

		dates = append(dates, date)
	}

	return dates
}

func TestStoreCreate(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	event, err := store.Create(t.Context(), r)
	is.NoErr(err)

	query := EventQuery{
		UserID: event.UserID,
		Period: PeriodDay,
	}

	events, err := store.Query(t.Context(), query)
	is.NoErr(err)

	is.Equal(len(events), 1)
	is.Equal(events[0], event)
}

func TestStoreUpdate(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	event, err := store.Create(t.Context(), r)
	is.NoErr(err)

	event.Text = "Text was updated..."

	updated, err := store.Update(t.Context(), event)
	is.NoErr(err)

	is.Equal(event, updated)
}

func TestStoreDelete(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	event, err := store.Create(t.Context(), CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	})
	is.NoErr(err)

	event, err = store.Delete(t.Context(), DeleteRequest{
		EventID: event.ID,
		UserID:  event.UserID,
	})
	is.NoErr(err)

	events, err := store.Query(t.Context(), EventQuery{
		UserID: event.UserID,
		Period: PeriodDay,
	})
	is.NoErr(err)

	is.Equal(len(events), 0)
}

func TestStoreQueryDay(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	event, err := store.Create(t.Context(), r)
	is.NoErr(err)

	events, err := store.Query(t.Context(), EventQuery{
		UserID: event.UserID,
		Period: PeriodDay,
	})
	is.NoErr(err)

	is.Equal(len(events), 1)
}

func TestStoreQueryEmptyDay(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now().Add(time.Hour * -24),
		Text:   "Hello, World!",
	}

	event, err := store.Create(t.Context(), r)
	is.NoErr(err)

	events, err := store.Query(t.Context(), EventQuery{
		UserID: event.UserID,
		Period: PeriodDay,
	})
	is.NoErr(err)

	is.Equal(len(events), 0)
}

func TestStoreInvalidQuery(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	_, err := store.Query(t.Context(), EventQuery{
		Period: 111,
	})
	is.Equal(err, ErrUserNotFound)
}

func TestStoreInvalidQuery2(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello",
	}

	_, err := store.Create(t.Context(), r)
	is.NoErr(err)

	_, err = store.Query(t.Context(), EventQuery{
		UserID: r.UserID,
		Period: -1,
	})

	is.Equal(err, ErrInvalidPeriod)
}

func TestStoreQueryWeek(t *testing.T) {
	var (
		is       = is.New(t)
		store    = NewInMemoryEventStore()
		userID   = uuid.New()
		weekDays = weekRange(time.Now())
		requests = make([]CreateRequest, 0, len(weekDays))
	)

	for i, day := range weekDays {
		requests = append(requests, CreateRequest{
			UserID: userID,
			Date:   day,
			Text:   fmt.Sprintf("Hello, World! – %d", i+1),
		})
	}

	var created []Event

	for _, req := range requests {
		event, err := store.Create(t.Context(), req)
		is.NoErr(err)

		created = append(created, event)
	}

	events, err := store.Query(t.Context(), EventQuery{
		UserID: userID,
		Period: PeriodWeek,
	})
	is.NoErr(err)

	is.Equal(len(events), len(weekDays))
	is.Equal(events, created)
}

func TestStoreQueryMonth(t *testing.T) {
	var (
		is        = is.New(t)
		store     = NewInMemoryEventStore()
		userID    = uuid.New()
		monthDays = monthRange(time.Now())
		requests  = make([]CreateRequest, 0, len(monthDays))
	)

	for i, day := range monthDays {
		requests = append(requests, CreateRequest{
			UserID: userID,
			Date:   day,
			Text:   fmt.Sprintf("Hello, World! – %d", i+1),
		})
	}

	var created []Event

	for _, req := range requests {
		event, err := store.Create(t.Context(), req)
		is.NoErr(err)

		created = append(created, event)
	}

	events, err := store.Query(t.Context(), EventQuery{
		UserID: userID,
		Period: PeriodMonth,
	})
	is.NoErr(err)

	is.Equal(len(events), len(monthDays))
	is.Equal(events, created)
}

func TestStoreQueryAll(t *testing.T) {
	var (
		is     = is.New(t)
		store  = NewInMemoryEventStore()
		userID = uuid.New()
	)

	var created []Event

	for n := range 100 {
		event, err := store.Create(t.Context(), CreateRequest{
			UserID: userID,
			Date:   time.Now().Add(time.Hour * 24 * time.Duration(n+1)),
		})
		is.NoErr(err)

		created = append(created, event)
	}

	events, err := store.Query(t.Context(), EventQuery{
		UserID: userID,
		Period: PeriodNone,
	})
	is.NoErr(err)

	is.Equal(len(events), 100)
	is.Equal(events, created)
}

func TestUpdateUserNotFound(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	_, err := store.Create(t.Context(), r)
	is.NoErr(err)

	_, err = store.Update(t.Context(), Event{
		UserID: uuid.Nil,
	})
	is.Equal(err, ErrUserNotFound)
}

func TestUpdateEventNotFound(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	_, err := store.Create(t.Context(), r)
	is.NoErr(err)

	_, err = store.Update(t.Context(), Event{
		UserID: r.UserID,
		ID:     uuid.Nil,
	})
	is.Equal(err, ErrEventNotFound)
}

func TestDeleteUserNotFound(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	_, err := store.Create(t.Context(), r)
	is.NoErr(err)

	_, err = store.Delete(t.Context(), DeleteRequest{
		UserID: uuid.Nil,
	})
	is.Equal(err, ErrUserNotFound)
}

func TestDeleteEventNotFound(t *testing.T) {
	is := is.New(t)

	store := NewInMemoryEventStore()

	r := CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now(),
		Text:   "Hello, World!",
	}

	_, err := store.Create(t.Context(), r)
	is.NoErr(err)

	_, err = store.Delete(t.Context(), DeleteRequest{
		EventID: uuid.Nil,
		UserID:  r.UserID,
	})
	is.Equal(err, ErrEventNotFound)
}
