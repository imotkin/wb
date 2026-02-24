package main

import (
	"log"

	"github.com/imotkin/L0/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start app: %v\n", err)
	}
}
