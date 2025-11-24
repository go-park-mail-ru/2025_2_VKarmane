package financeservice

import (
	"log/slog"
	"net"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	fin "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/grpc"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	finrepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/repository"
	finsvc "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/service"
	finusecase "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/usecase"
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

	lis, err := net.Listen("tcp", ":8110")
	if err != nil {
		appLogger.Error("failed to start FinanceService", "error", err)
		return err
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			interceptors.LoggerInterceptor(appLogger),
		),
	)

	grpc_prometheus.Register(srv)
	grpc_prometheus.EnableHandlingTimeHistogram()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		appLogger.Info("Metrics server started on :8910")
		if err := http.ListenAndServe(":8910", mux); err != nil {
			appLogger.Error("Metrics server failed", err)
		}
	}()

	db, err := finrepo.NewDBConnection(config.GetDatabaseDSN())
	if err != nil {
		appLogger.Error("FinanceService failed to connect to DB", "error", err)
		return err
	}
	store := finrepo.NewPostgresRepository(db)
	svc := finsvc.NewService(store, clock)

	uc := finusecase.NewFinanceUseCase(svc)
	financeService := fin.NewFinanceServer(uc)

	finpb.RegisterFinanceServiceServer(srv, financeService)

	appLogger.Info("FinanceService started on :8110")
	srv.Serve(lis)

	return nil
}
