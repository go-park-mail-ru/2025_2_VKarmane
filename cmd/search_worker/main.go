package main

import (
	"log"

	srch "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker"
)

func main() {
	if err := srch.Run(); err != nil {
		log.Fatalf("AuthService failed to run: %v", err)
	}
}
