package router

import (
	"net/http"

	"github.com/imotkin/L0/internal/api/handler"
	"github.com/imotkin/L0/internal/metrics"
)

func New(h *handler.Handler, templatePath string) *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("GET /order/{id}", h.GetOrder())
	r.Handle("GET /orders", h.GetList())
	r.Handle("GET /search", h.IndexPage(templatePath))
	r.Handle("/metrics", metrics.Handler())

	return r
}
