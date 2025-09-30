package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
