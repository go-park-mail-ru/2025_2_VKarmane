package authservice

import (
	"net"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/impl"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	authsvc "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/service"
	repo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/repository"
	authusecase "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/interceptors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
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
		appLogger.Error("failed to start AuthService %w", err)
		return err
	}


	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.LoggerInterceptor(appLogger)))
	db, err := repo.NewDBConnection(config.GetDatabaseDSN())
	if err != nil {
		appLogger.Error("AuthService failed to connect do DB %w", err)
		return err
	}
	store := repo.NewPostgresRepository(db)
	svc := authsvc.NewService(store, config.JWTSecret, clock)
	

	
	uc := authusecase.NewAuthUseCase(svc, config.JWTSecret, clock)
	authService := auth.NewAuthServer(uc)

	authpb.RegisterAuthServiceServer(srv, authService)


	srv.Serve(lis)
	
	return nil
}