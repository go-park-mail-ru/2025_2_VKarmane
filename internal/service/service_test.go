package service

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestNewService_WiresServices(t *testing.T) {
	store, err := repository.NewStore()
	require.NoError(t, err)
	s := NewService(store)
	require.NotNil(t, s)
	require.NotNil(t, s.AuthUC)
	require.NotNil(t, s.BalanceUC)
	require.NotNil(t, s.BudgetUC)
}
