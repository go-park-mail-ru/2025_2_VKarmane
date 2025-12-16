package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service"
)



func main() {
	if err := budgservice.Run(); err != nil {
		log.Fatalf("BudgetService failed to run: %v", err)
	}
}

