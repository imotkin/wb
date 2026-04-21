package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/imotkin/L2/18/internal/calendar"
)

type Handler struct {
	log     *slog.Logger
	service calendar.Service
}

func NewHandler(log *slog.Logger, service calendar.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) json(w http.ResponseWriter, v any, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.log.Error("encode json response", "err", err)
	}

	switch value := v.(type) {
	case *ErrorResponse:
		h.log.Error("failed to handle request", "err", value.Error)
	}
}

func (h *Handler) decodeEvent(r *http.Request) (calendar.Event, error) {
	var event calendar.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return calendar.Event{}, err
	}

	err = event.Validate()
	if err != nil {
		return calendar.Event{}, err
	}

	return event, nil
}

func (h *Handler) getUser(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.URL.Query().Get("user_id"))
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req calendar.CreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.json(w, NewError("invalid delete request"), http.StatusBadRequest)
		return
	}

	event, err := h.service.CreateEvent(context.Background(), req)
	if err != nil {
		h.json(w, NewError("failed to add event"), http.StatusBadRequest)
		return
	}

	h.json(w, event, http.StatusOK)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := h.decodeEvent(r)
	if err != nil {
		h.json(w, NewError("invalid event body"), http.StatusBadRequest)
		return
	}

	event, err = h.service.UpdateEvent(r.Context(), event)
	if err != nil {
		h.json(w, NewError("failed to update event"), http.StatusBadRequest)
		return
	}

	h.json(w, event, http.StatusOK)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req calendar.DeleteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.json(w, NewError("invalid delete request"), http.StatusBadRequest)
		return
	}

	event, err := h.service.DeleteEvent(r.Context(), req)
	if err != nil {
		h.json(w, NewError("failed to delete event"), http.StatusBadRequest)
		return
	}

	h.json(w, event, http.StatusOK)
}

func (h *Handler) DayEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUser(r)
	if err != nil {
		h.json(w, NewError("invalid user id"), http.StatusBadRequest)
		return
	}

	events, err := h.service.DayEvents(r.Context(), userID)
	if err != nil {
		h.json(w, NewError("failed to get events"), http.StatusBadRequest)
		return
	}

	h.json(w, events, http.StatusOK)
}

func (h *Handler) WeekEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUser(r)
	if err != nil {
		h.json(w, NewError("invalid user id"), http.StatusBadRequest)
		return
	}

	events, err := h.service.WeekEvents(r.Context(), userID)
	if err != nil {
		h.json(w, NewError("failed to get events"), http.StatusBadRequest)
		return
	}

	h.json(w, events, http.StatusOK)
}

func (h *Handler) MonthEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUser(r)
	if err != nil {
		h.json(w, NewError("invalid user id"), http.StatusBadRequest)
		return
	}

	events, err := h.service.MonthEvents(r.Context(), userID)
	if err != nil {
		h.json(w, NewError("failed to get events"), http.StatusBadRequest)
		return
	}

	h.json(w, events, http.StatusOK)
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUser(r)
	if err != nil {
		h.json(w, NewError("invalid user id"), http.StatusBadRequest)
		return
	}

	events, err := h.service.Events(r.Context(), userID)
	if err != nil {
		h.json(w, NewError("failed to get events"), http.StatusBadRequest)
		return
	}

	h.json(w, events, http.StatusOK)
}
