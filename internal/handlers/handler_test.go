package handlers

import (
	"context"
	"io"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	budgetpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	imagerepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository"
	imageservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
)

type testImageStorage struct{}

func (testImageStorage) UploadImage(context.Context, io.Reader, string, int64, string) (string, error) {
	return "id", nil
}

func (testImageStorage) GetImageURL(context.Context, string) (string, error) {
	return "url", nil
}

func (testImageStorage) DeleteImage(context.Context, string) error {
	return nil
}

func (testImageStorage) ImageExists(context.Context, string) (bool, error) {
	return true, nil
}

func newTestUseCase() *usecase.UseCase {
	imageService := imageservice.NewService(testImageStorage{})
	svc := &service.Service{ImageUC: imageService}
	return usecase.NewUseCase(svc, "secret")
}

func closeLogger(t *testing.T, l logger.Logger) {
	t.Helper()
	if closable, ok := l.(*logger.SlogLogger); ok {
		require.NoError(t, closable.Close())
	}
}

func TestNewHandlerInitializesDependencies(t *testing.T) {
	uc := newTestUseCase()
	l := logger.NewSlogLogger()
	defer closeLogger(t, l)

	h := NewHandler(uc, l, nil, nil, nil, nil)
	require.NotNil(t, h.balanceHandler)
	require.NotNil(t, h.budgetHandler)
	require.NotNil(t, h.authHandler)
	require.NotNil(t, h.opHandler)
	require.NotNil(t, h.categoryHandler)
	require.NotNil(t, h.profileHandler)
	require.NotNil(t, h.registrator)
}

func TestHandlerRegisterAddsRoutes(t *testing.T) {
	uc := newTestUseCase()
	l := logger.NewSlogLogger()
	defer closeLogger(t, l)

	h := NewHandler(uc, l, nil, nil, nil, nil)
	publicRouter := mux.NewRouter()
	protectedRouter := mux.NewRouter()

	h.Register(publicRouter, protectedRouter, nil, nil, nil, nil)

	publicCount := countRoutes(publicRouter)
	protectedCount := countRoutes(protectedRouter)

	require.Greater(t, publicCount, 0)
	require.Greater(t, protectedCount, 0)
}

func TestRegistratorRegisterAll(t *testing.T) {
	uc := newTestUseCase()
	l := logger.NewSlogLogger()
	defer closeLogger(t, l)

	reg := NewRegistrator(uc, l)
	publicRouter := mux.NewRouter()
	protectedRouter := mux.NewRouter()

	var dummyAuth authpb.AuthServiceClient
	var dummyBudget budgetpb.BudgetServiceClient
	var dummyFin finpb.FinanceServiceClient

	reg.RegisterAll(publicRouter, protectedRouter, uc, l, dummyAuth, dummyBudget, dummyFin, nil)

	require.Greater(t, countRoutes(publicRouter), 0)
	require.Greater(t, countRoutes(protectedRouter), 0)
}

func countRoutes(r *mux.Router) int {
	count := 0
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		count++
		return nil
	})
	return count
}

var _ imagerepo.ImageStorage = (*testImageStorage)(nil)