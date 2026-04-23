package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

const consoleHint = "> "

type Client struct {
	addr    string
	conn    net.Conn
	cancel  context.CancelFunc
	timeout time.Duration
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Close() {
	if c.conn == nil {
		return
	}

	c.conn.Close()
}

func (c *Client) ReadInput(ctx context.Context) <-chan string {
	in := make(chan string)

	go func() {
		defer func() {
			close(in)
			c.cancel()
		}()

		fmt.Print(consoleHint)

		scan := bufio.NewScanner(os.Stdin)

		for scan.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			text := scan.Text()

			fmt.Printf("[READ] -> %q\n", text)

			in <- text
		}
	}()

	return in
}

func (c *Client) Run(ctx context.Context, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", c.addr, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect server: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	c.conn = conn
	c.cancel = cancel
	c.timeout = timeout

	in := c.ReadInput(ctx)

	go c.Listen(ctx)
	go c.Send(ctx, in)

	<-ctx.Done()
	c.Close()

	return nil
}

func (c *Client) Listen(ctx context.Context) {
	defer c.cancel()

	scan := bufio.NewScanner(c.conn)

	for scan.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		text := scan.Text()

		fmt.Printf("[RECV] <- %q\n%s", text, consoleHint)
	}
}

func (c *Client) Send(ctx context.Context, in <-chan string) {
	defer c.cancel()

	for {
		select {
		case s, ok := <-in:
			if !ok {
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(c.timeout))

			_, err := fmt.Fprintln(c.conn, s)
			if err != nil {
				return
			}

			fmt.Printf("[SEND] -> %d bytes\n", len(s))
		case <-ctx.Done():
			return
		}
	}
}
