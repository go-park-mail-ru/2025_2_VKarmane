package budgservice

import (
	"log/slog"
	"net"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	bdg "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/grpc"
	"github.com/go-park-mail-ru/2025_2_VKarmane/pkg/interceptors"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	bdgrepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/repository"
	bdgsvc "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/service"
	bdgusecase "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
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
		appLogger.Info("Metrics server started on :10200")
		if err := http.ListenAndServe(":8800", mux); err != nil {
			appLogger.Error("Metrics server failed", err)
		}
	}()

	db, err := bdgrepo.NewDBConnection(config.GetDatabaseDSN())
	if err != nil {
		appLogger.Error("BudgetService failed to connect to DB %w", err)
		return err
	}
	store := bdgrepo.NewPostgresRepository(db)
	svc := bdgsvc.NewService(store, clock)

	uc := bdgusecase.NewBudgetUseCase(svc)
	bdgService := bdg.NewBudgetServer(uc)

	bdgpb.RegisterBudgetServiceServer(srv, bdgService)

	srv.Serve(lis)

	return nil
}
