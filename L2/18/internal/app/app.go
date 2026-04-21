package app

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/imotkin/L2/18/internal/calendar"
	"github.com/imotkin/L2/18/internal/config"
	"github.com/imotkin/L2/18/internal/http"
)

var configPath = flag.String("config", "config.example.yaml", "path to config file")

func Run() error {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Parse(*configPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	var (
		log   = slog.Default()
		store = calendar.NewInMemoryEventStore()
		svc   = calendar.NewService(store)
		h     = http.NewHandler(log, svc)
		r     = http.NewRouter(log, h)
		srv   = http.NewServer(log, cfg.Server, r)
	)

	err = srv.Start(ctx)
	if err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}
