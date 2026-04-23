package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	defaultTimeout = time.Second * 10
)

var (
	flagHost    = flag.String("host", "", "The host of TCP server")
	flagPort    = flag.String("port", "", "The port of TCP server")
	flagTimeout = flag.Duration("timeout", defaultTimeout, "The timeout of server response")
)

func run() error {
	flag.Parse()

	if *flagHost == "" || *flagPort == "" {
		return errors.New("invalid connection args")
	}

	addr := net.JoinHostPort(*flagHost, *flagPort)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := NewClient(addr)

	return client.Run(ctx, *flagTimeout)
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("telnet: %v\n", err)
	}
}
