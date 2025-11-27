package main

import (
	"fmt"
	"log"

	srch "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker"
)

func main() {
	if err := srch.Run(); err != nil {
		log.Fatal(fmt.Sprintf("AuthService failed to run: %v", err))
	}
}
