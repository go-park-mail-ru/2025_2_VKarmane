package budgservice

import (
	"net"
	"log/slog"

	"google.golang.org/grpc"

	bdg "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/impl"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	bdgsvc "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/service"
	store "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/store"
	bdgusecase "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/interceptors"
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


	lis, err := net.Listen("tcp", ":8100")
	if err != nil {
		appLogger.Error("failed to start BudgetService %w", err)
		return err
	}


	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.LoggerInterceptor(appLogger)))
	
	store, err := store.NewPostgresStore(config.GetDatabaseDSN())
	svc := bdgsvc.NewService(store, clock)
	

	
	uc := bdgusecase.NewBudgetUseCase(svc)
	bdgService := bdg.NewBudgetServer(uc)

	bdgpb.RegisterBudgetServiceServer(srv, bdgService)


	srv.Serve(lis)
	
	return nil
}