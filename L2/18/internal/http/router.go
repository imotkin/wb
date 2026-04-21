package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(log *slog.Logger, h *Handler) http.Handler {
	m := chi.NewMux()

	m.Use(LoggingMiddleware(log))

	m.Post("/create_event", h.CreateEvent)
	m.Post("/update_event", h.UpdateEvent)
	m.Post("/delete_event", h.DeleteEvent)

	m.Get("/events_for_day", h.DayEvents)
	m.Get("/events_for_week", h.WeekEvents)
	m.Get("/events_for_month", h.MonthEvents)

	m.Get("/events", h.Events)

	return m
}
