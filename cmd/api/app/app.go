package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"

	"github.com/gorilla/mux"
)

func Run() {
	// TODO: use environment variables
	jwtSecret := "your-secret-key"

	appLogger, err := logger.NewSlogLoggerWithFileAndConsole("logs/app.log", slog.LevelInfo)
	if err != nil {
		appLogger = logger.NewSlogLogger()
	}

	store := repository.NewStore()
	service := service.NewService(store)
	usecase := usecase.NewUseCase(service, store, jwtSecret)
	handler := handlers.NewHandler(usecase, appLogger)

	r := mux.NewRouter()

	r.Use(middleware.LoggerMiddleware(appLogger))

	r.Use(middleware.RequestLoggerMiddleware(appLogger))

	r.Use(middleware.SecurityLoggerMiddleware(appLogger))

	public := r.PathPrefix("/api/v1").Subrouter()
	public.HandleFunc("/auth/register", handler.Register).Methods(http.MethodPost)
	public.HandleFunc("/auth/login", handler.Login).Methods(http.MethodPost)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	protected.HandleFunc("/auth/logout", handler.Logout).Methods(http.MethodPost)
	protected.HandleFunc("/profile", handler.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/budgets", handler.GetListBudgets).Methods(http.MethodGet)
	protected.HandleFunc("/budget/{id}", handler.GetBudgetByID).Methods(http.MethodGet)
	protected.HandleFunc("/balance", handler.GetListBalance).Methods(http.MethodGet)
	protected.HandleFunc("/balance/{id}", handler.GetBalanceByAccountID).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Server running at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
