package main

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service"
)



func main() {
	if err := budgservice.Run(); err != nil {
		log.Fatal(fmt.Sprintf("BudgetService failed to run: %w", err))
	}
}

