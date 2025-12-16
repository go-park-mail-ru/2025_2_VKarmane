package main

import (
	"log"

	budgservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service"
)

func main() {
	if err := budgservice.Run(); err != nil {
		log.Fatalf("FinanceService failed to run: %v", err)
	}
}
