package budget

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_GetBudgets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := NewBudgetUseCase(mockSvc)
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	userID := 1
	mockResp := &bdgpb.ListBudgetsResponse{
		Budgets: []*bdgpb.Budget{
			{Id: 1, UserId: 1, Sum: 1000},
			{Id: 2, UserId: 1, Sum: 500},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().GetBudgets(gomock.Any(), userID).Return(mockResp, nil)
		resp, err := uc.GetBudgets(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, mockResp, resp)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.EXPECT().GetBudgets(gomock.Any(), userID).Return(nil, errors.New("db fail"))
		resp, err := uc.GetBudgets(ctx, userID)
		assert.Nil(t, resp)
		assert.ErrorContains(t, err, "budget.GetBudgetsForUser")
	})
}

func TestUseCase_GetBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := NewBudgetUseCase(mockSvc)
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	budgetID := 1
	userID := 1
	mockBudget := &bdgpb.Budget{Id: int32(budgetID), UserId: int32(userID), Sum: 1000}

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().GetBudgetByID(gomock.Any(), budgetID, userID).Return(mockBudget, nil)
		b, err := uc.GetBudget(ctx, budgetID, userID)
		assert.NoError(t, err)
		assert.Equal(t, mockBudget, b)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.EXPECT().GetBudgetByID(gomock.Any(), budgetID, userID).Return(nil, errors.New("db fail"))
		b, err := uc.GetBudget(ctx, budgetID, userID)
		assert.Nil(t, b)
		assert.ErrorContains(t, err, "budget.GetBudgetByID")
	})
}

func TestUseCase_CreateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := NewBudgetUseCase(mockSvc)
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	userID := 1
	req := models.CreateBudgetRequest{CategoryID: 1, Amount: 100}
	mockBudget := &bdgpb.Budget{Id: 10, UserId: int32(userID), Sum: 100}

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().CreateBudget(gomock.Any(), req, userID).Return(mockBudget, nil)
		b, err := uc.CreateBudget(ctx, req, userID)
		assert.NoError(t, err)
		assert.Equal(t, mockBudget, b)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.EXPECT().CreateBudget(gomock.Any(), req, userID).Return(nil, errors.New("active budget exists"))
		b, err := uc.CreateBudget(ctx, req, userID)
		assert.Nil(t, b)
		assert.ErrorContains(t, err, "budget.CreateBudget")
	})
}

func TestUseCase_UpdateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := NewBudgetUseCase(mockSvc)
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	req := models.UpdatedBudgetRequest{UserID: 1, BudgetID: 1}
	mockBudget := &bdgpb.Budget{Id: 1, UserId: 1, Sum: 200}

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().UpdateBudget(gomock.Any(), req).Return(mockBudget, nil)
		b, err := uc.UpdateBudget(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, mockBudget, b)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.EXPECT().UpdateBudget(gomock.Any(), req).Return(nil, errors.New("db fail"))
		b, err := uc.UpdateBudget(ctx, req)
		assert.Nil(t, b)
		assert.ErrorContains(t, err, "budget.UpdateBudget")
	})
}

func TestUseCase_DeleteBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := NewBudgetUseCase(mockSvc)
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	budgetID := 1
	userID := 1
	mockBudget := &bdgpb.Budget{Id: int32(budgetID), UserId: int32(userID), Sum: 1000}

	t.Run("success", func(t *testing.T) {
		mockSvc.EXPECT().DeleteBudget(gomock.Any(), budgetID, userID).Return(mockBudget, nil)
		b, err := uc.DeleteBudget(ctx, budgetID, userID)
		assert.NoError(t, err)
		assert.Equal(t, mockBudget, b)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.EXPECT().DeleteBudget(gomock.Any(), budgetID, userID).Return(nil, errors.New("db fail"))
		b, err := uc.DeleteBudget(ctx, budgetID, userID)
		assert.Nil(t, b)
		assert.ErrorContains(t, err, "budget.DeleteBudget")
	})
}
