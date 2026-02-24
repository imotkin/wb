package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

func (h *Handler) error(w http.ResponseWriter, msg string, code int, err error) {
	h.log.Error(err, msg)
	h.response(w, ErrorMessage{
		Message:       msg,
		StatusCode:    code,
		StatusMessage: http.StatusText(code),
	}, code)
}

func (h *Handler) response(w http.ResponseWriter, v any, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.log.Error(err, "failed to send json response")
	}
}
