package app

import (
	"context"
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
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() error {
	config := LoadConfig()

	appLogger, err := logger.NewSlogLoggerWithFileAndConsole("logs/app.log", slog.LevelInfo)
	if err != nil {
		appLogger = logger.NewSlogLogger()
	}

	store, err := repository.NewPostgresStore(config.GetDatabaseDSN())
	if err != nil {
		return err
	}
	defer func() {
		if err := store.Close(); err != nil {
			appLogger.Error("Failed to close database connection", "error", err)
		}
	}()
	service := service.NewService(store, config.JWTSecret)
	usecase := usecase.NewUseCase(service, store, config.JWTSecret)
	handler := handlers.NewHandler(usecase, appLogger)

	r := mux.NewRouter()

	// CORS middleware должен быть первым
	corsOrigins := config.GetCORSOrigins()
	appLogger.Info("CORS Configuration", "origins", corsOrigins, "is_production", config.IsProduction())
	r.Use(middleware.CORSMiddleware(corsOrigins, appLogger))

	r.Use(middleware.LoggerMiddleware(appLogger))

	r.Use(middleware.RequestLoggerMiddleware(appLogger))

	r.Use(middleware.SecurityLoggerMiddleware(appLogger))

	public := r.PathPrefix("/api/v1").Subrouter()

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.CORSMiddleware(corsOrigins, appLogger))
	protected.Use(middleware.LoggerMiddleware(appLogger))
	protected.Use(middleware.RequestLoggerMiddleware(appLogger))
	protected.Use(middleware.SecurityLoggerMiddleware(appLogger))
	protected.Use(middleware.AuthMiddleware(config.JWTSecret))
	handler.Register(public, protected)

	// Swagger документация
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Добавляем обработку OPTIONS запросов для всех маршрутов (для preflight запросов)
	r.PathPrefix("/api/v1").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			appLogger.Info("Handling OPTIONS request", "path", r.URL.Path)
			w.WriteHeader(http.StatusOK)

			return
		}
		http.NotFound(w, r)
	}).Methods(http.MethodOptions)
	srv := &http.Server{
		Addr:    config.GetServerAddress(),
		Handler: r,
	}

	go func() {
		appLogger.Info("Server running at", "address", config.GetServerAddress(), "https_enabled", config.HTTPS.Enabled)

		var err error
		if config.HTTPS.Enabled {
			appLogger.Info("Starting HTTPS server", "cert_file", config.HTTPS.CertFile, "key_file", config.HTTPS.KeyFile)
			err = srv.ListenAndServeTLS(config.HTTPS.CertFile, config.HTTPS.KeyFile)
		} else {
			appLogger.Info("Starting HTTP server")
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited")
	return nil
}
