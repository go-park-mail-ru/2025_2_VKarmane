package main

import (
	"fmt"
	"log"

	budgservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service"
)

func main() {
	if err := budgservice.Run(); err != nil {
		log.Fatal(fmt.Sprintf("FinanceService failed to run: %v", err))
	}
}



