package authservice

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	server "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/grpc"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	repo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/repository"
	authusecase "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/pkg/interceptors"
)

func Run() error {
	config := config.LoadConfig()
	clock := clock.RealClock{}

	appLogger, err := logger.NewSlogLoggerWithFileAndConsole("logs/app.log", slog.LevelInfo)
	if err != nil {
		appLogger = logger.NewSlogLogger()
	}

	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		appLogger.Error("failed to start AuthService ", "error", err)
		return err
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			interceptors.LoggerInterceptor(appLogger),
		),
		// grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	grpc_prometheus.Register(srv)
	grpc_prometheus.EnableHandlingTimeHistogram()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		appLogger.Info("Metrics server started on :10300")
		if err := http.ListenAndServe(":8700", mux); err != nil {
			appLogger.Error("Metrics server failed", err)
		}
	}()

	db, err := repo.NewDBConnection(config.GetDatabaseDSN())
	if err != nil {
		appLogger.Error("AuthService failed to connect do DB ", "errror", err)
		return err
	}
	store := repo.NewPostgresRepository(db)

	uc := authusecase.NewAuthUseCase(store, config.JWTSecret, clock)
	authService := server.NewAuthServer(uc)

	authpb.RegisterAuthServiceServer(srv, authService)

	srv.Serve(lis)

	return nil
}
