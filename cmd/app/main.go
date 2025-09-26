package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// TODO: use environment variables
	jwtSecret := "your-secret-key"

	store := repository.NewStore()
	service := service.NewService(store)
	usecase := usecase.NewUseCase(service, store, jwtSecret)
	handler := handlers.NewHandler(usecase)

	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/api/v1/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", handler.Login).Methods("POST")

	// Защищенные маршруты
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	protected.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	protected.HandleFunc("/budgets", handler.GetListBudgets).Methods("GET")
	protected.HandleFunc("/budget/{id}", handler.GetBudgetByID).Methods("GET")
	protected.HandleFunc("/balance", handler.GetListBalance).Methods("GET")
	protected.HandleFunc("/balance/{id}", handler.GetBalanceByAccountID).Methods("GET")

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
