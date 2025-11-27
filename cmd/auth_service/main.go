package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service"
)



func main() {
	if err := authservice.Run(); err != nil {
		log.Fatalf("AuthService failed to run: %v", err)
	}
}

