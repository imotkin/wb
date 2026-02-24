package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/imotkin/L0/internal/logger"
)

type Server struct {
	srv *http.Server
	log logger.Logger
}

func New(log logger.Logger, h http.Handler, host, port string) *Server {
	addr := net.JoinHostPort(host, port)

	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      h,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			IdleTimeout:  time.Second * 60,
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
			return err
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
