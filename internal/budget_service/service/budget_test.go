package budget

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/errors"
	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetBudgets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBudgetRepository(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := 1

	budgets := []bdgmodels.Budget{
		{ID: 1, UserID: userID, Amount: 100, CurrencyID: 1, PeriodStart: time.Now(), PeriodEnd: time.Now().Add(24 * time.Hour)},
		{ID: 2, UserID: userID, Amount: 200, CurrencyID: 1, PeriodStart: time.Now(), PeriodEnd: time.Now().Add(24 * time.Hour)},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return(budgets, nil)

		resp, err := svc.GetBudgets(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Budgets))
		assert.Equal(t, int32(1), resp.Budgets[0].Id)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return(nil, errors.New("db error"))

		resp, err := svc.GetBudgets(ctx, userID)
		assert.Nil(t, resp)
		assert.ErrorContains(t, err, "Failed to get budgets for user")
	})
}

func TestService_GetBudgetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBudgetRepository(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := 1
	budgetID := 1

	budget := bdgmodels.Budget{ID: budgetID, UserID: userID, Amount: 100}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{budget}, nil).Times(2)

		resp, err := svc.GetBudgetByID(ctx, budgetID, userID)
		assert.NoError(t, err)
		assert.Equal(t, int32(budgetID), resp.Id)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{}, nil)

		resp, err := svc.GetBudgetByID(ctx, budgetID, userID)
		assert.Nil(t, resp)
		assert.Equal(t, bdgerrors.ErrForbidden, err)
	})

	t.Run("repo_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockBudgetRepository(ctrl)
		s := NewService(mockRepo, clock.RealClock{})

		ctx := context.Background()
		userID := 1
		budgetID := 1

		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{
			{ID: budgetID, UserID: userID},
		}, nil).Times(1)

		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return(nil, errors.New("some db error")).Times(1)

		_, err := s.GetBudgetByID(ctx, budgetID, userID)
		if err == nil || !strings.Contains(err.Error(), "budget.GetBudgetByID") {
			t.Fatalf("expected error containing 'budget.GetBudgetByID', got %v", err)
		}
	})

}

func TestService_CreateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBudgetRepository(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := 1
	req := bdgmodels.CreateBudgetRequest{CategoryID: 1, Amount: 100}
	budget := bdgmodels.Budget{ID: 1, UserID: userID, Amount: 100}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().CreateBudget(ctx, gomock.Any()).Return(budget, nil)

		resp, err := svc.CreateBudget(ctx, req, userID)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), resp.Id)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().CreateBudget(ctx, gomock.Any()).Return(bdgmodels.Budget{}, errors.New("db fail"))

		resp, err := svc.CreateBudget(ctx, req, userID)
		assert.Nil(t, resp)
		assert.ErrorContains(t, err, "Failed to create budget")
	})
}

func TestService_UpdateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBudgetRepository(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := 1
	budgetID := 1
	amount := 200.0
	req := bdgmodels.UpdatedBudgetRequest{UserID: userID, BudgetID: budgetID, Amount: &amount}
	updatedBudget := bdgmodels.Budget{ID: budgetID, UserID: userID, Amount: 200}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{{ID: budgetID, UserID: userID}}, nil)
		mockRepo.EXPECT().UpdateBudget(ctx, req).Return(updatedBudget, nil)

		resp, err := svc.UpdateBudget(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int32(budgetID), resp.Id)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{}, nil)

		resp, err := svc.UpdateBudget(ctx, req)
		assert.Nil(t, resp)
		assert.Equal(t, bdgerrors.ErrForbidden, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{{ID: budgetID, UserID: userID}}, nil)
		mockRepo.EXPECT().UpdateBudget(ctx, req).Return(bdgmodels.Budget{}, errors.New("db fail"))

		resp, err := svc.UpdateBudget(ctx, req)
		assert.Nil(t, resp)
		assert.ErrorContains(t, err, "Failed to update budget")
	})
}

func TestService_DeleteBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBudgetRepository(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := 1
	budgetID := 1
	deletedBudget := bdgmodels.Budget{ID: budgetID, UserID: userID, Amount: 100}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{{ID: budgetID, UserID: userID}}, nil)
		mockRepo.EXPECT().DeleteBudget(ctx, budgetID).Return(deletedBudget, nil)

		resp, err := svc.DeleteBudget(ctx, userID, budgetID)
		assert.NoError(t, err)
		assert.Equal(t, int32(budgetID), resp.Id)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{}, nil)

		resp, err := svc.DeleteBudget(ctx, userID, budgetID)
		assert.Nil(t, resp)
		assert.Equal(t, bdgerrors.ErrForbidden, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetBudgetsByUser(ctx, userID).Return([]bdgmodels.Budget{{ID: budgetID, UserID: userID}}, nil)
		mockRepo.EXPECT().DeleteBudget(ctx, budgetID).Return(bdgmodels.Budget{}, errors.New("db fail"))

		resp, err := svc.DeleteBudget(ctx, userID, budgetID)
		assert.Nil(t, resp)
		assert.ErrorContains(t, err, "Failed to delete budget")
	})
}
