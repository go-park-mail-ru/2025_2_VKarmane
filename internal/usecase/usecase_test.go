package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/stretchr/testify/require"
)

func TestNewUseCase_WiresDependencies(t *testing.T) {
	store, err := repository.NewStore()
	require.NoError(t, err)
	svc := service.NewService(store)
	uc := NewUseCase(svc, store, "secret")
	require.NotNil(t, uc)
	require.NotNil(t, uc.AuthUC)
	require.NotNil(t, uc.BalanceUC)
	require.NotNil(t, uc.BudgetUC)
}
