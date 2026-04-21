package http

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(
				"got http request",
				slog.Group(
					"request",
					"method", r.Method,
					"url", r.URL.String(),
				),
			)
			next.ServeHTTP(w, r)
		})
	}
}
