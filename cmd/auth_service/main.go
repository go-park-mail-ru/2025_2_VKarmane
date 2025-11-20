package main

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service"
)



func main() {
	if err := authservice.Run(); err != nil {
		log.Fatal(fmt.Sprintf("AuthService failed to run: %v", err))
	}
}

