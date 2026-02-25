package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"

	"github.com/imotkin/L0/internal/entity"
	"github.com/imotkin/L0/internal/logger"
	"github.com/imotkin/L0/internal/metrics"
	"github.com/imotkin/L0/internal/service"
)

type Handler struct {
	s   service.Service
	log logger.Logger
}

func New(log logger.Logger, s service.Service) *Handler {
	return &Handler{s: s, log: log.With("source", "handler")}
}

func (h *Handler) GetOrder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequests()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			h.error(w, "invalid order id", http.StatusBadRequest, err)
			return
		}

		h.log.Info("got a new request", "id", id)

		order, err := h.s.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, entity.ErrOrderNotFound) {
				msg := fmt.Sprintf("order %q is not found", id)
				h.error(w, msg, http.StatusNotFound, err)
				return
			}

			h.error(w, "failed to get order", http.StatusInternalServerError, err)
			return
		}

		h.response(w, order, http.StatusOK)
	})
}

func (h *Handler) GetList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequests()

		orders, err := h.s.List(r.Context())
		if err != nil {
			h.error(w, "failed to get orders list", http.StatusInternalServerError, err)
			return
		}

		h.response(w, orders, http.StatusOK)
	})
}

func (h *Handler) IndexPage(templatePath string) http.Handler {
	tmpl := template.Must(template.ParseFiles(templatePath))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			h.log.Error(err, "failed to execute template")
		}
	})
}
