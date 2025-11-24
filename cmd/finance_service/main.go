package main

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service"
)

func main() {
	if err := financeservice.Run(); err != nil {
		log.Fatal(fmt.Sprintf("FinanceService failed to run: %v", err))
	}
}

