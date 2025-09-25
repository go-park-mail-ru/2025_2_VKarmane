package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	st := repository.NewStore()
	svc := service.NewService(st)
	h := handlers.NewHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/budgets", h.GetListBudgets).Methods("GET")
	r.HandleFunc("/api/v1/budget/{id}", h.GetBudgetByID).Methods("GET")
	r.HandleFunc("/api/v1/balance", h.GetListBalance).Methods("GET")
	r.HandleFunc("/api/v1/balance/{id}", h.GetBalanceByAccountID).Methods("GET")

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
