package main

import (
	"log"

	"github.com/imotkin/L2/18/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
