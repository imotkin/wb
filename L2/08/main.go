package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

var serverAddr = flag.String("server", "0.ru.pool.ntp.org", "The NTP server address")

func getTime(addr string) (time.Time, error) {
	response, err := ntp.Query(addr)
	if err != nil {
		return time.Time{}, fmt.Errorf("query time: %w", err)
	}

	return time.Now().Add(response.ClockOffset), nil
}

func main() {
	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags)

	currentTime, err := getTime(*serverAddr)
	if err != nil {
		logger.Fatalf("Failed to get current time: %s\n", err)
	}

	fmt.Printf("Current time -> %s\n", currentTime)
}
