package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	log *slog.Logger
}

func NewServer(log *slog.Logger, cfg *ServerConfig, h http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         cfg.Addr(),
			Handler:      h,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		log: log.With("source", "server"),
	}
}

func (s *Server) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.log.Info("started http server", "addr", s.srv.Addr)

		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Info("server", "error", err)
			return fmt.Errorf("start server listening: %w", err)
		}

		return nil
	})

	<-ctx.Done()

	s.log.Info("got an interrupt signal")

	g.Go(func() error {
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Second*5,
		)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	err := g.Wait()
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
