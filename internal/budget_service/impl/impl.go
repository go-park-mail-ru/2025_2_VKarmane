package budget

import (
	"context"
	// "errors"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	// "google.golang.org/protobuf/types/known/emptypb"

	budgetpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	// svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/errors"
	// "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetServerImpl struct {
	bdgUC BudgetUseCase
	budgetpb.UnimplementedBudgetServiceServer
}

func NewBudgetServer(bdgUC BudgetUseCase) *BudgetServerImpl {
	return &BudgetServerImpl{
		bdgUC: bdgUC,
	}
}


func (s *BudgetServerImpl) CreateBudget(ctx context.Context, req *budgetpb.CreateBudgetRequest) (*budgetpb.Budget, error) {
	budget, UserID := ProtoCreateRequestToModel(req)
	createdBdg, err := s.bdgUC.CreateBudget(ctx, budget, UserID)
	if err != nil {
		return nil, err
	}
	return createdBdg, nil
}

func (s *BudgetServerImpl) GetBudgetByID(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, UserID := ProtoBudgetReqToInts(req)
	bdg, err := s.bdgUC.GetBudget(ctx, budgetID, UserID)
	if err != nil {
		return nil, err
	}
	return bdg, nil
}

func (s *BudgetServerImpl) GetListBudgets(ctx context.Context, req *budgetpb.UserID) (*budgetpb.ListBudgetsResponse, error) {
	UserID := ProtoIDToInt(req)
	bdg, err := s.bdgUC.GetBudgets(ctx, UserID)
	if err != nil {
		return nil, err
	}
	return bdg, nil
}

func (s *BudgetServerImpl) UpdateBudget(ctx context.Context, req *budgetpb.UpdateBudgetRequest) (*budgetpb.Budget, error) {
	budget := ProtoUpdateRequestToModel(req)
	updatedBdg, err := s.bdgUC.UpdateBudget(ctx, budget)
	if err != nil {
		return nil, err
	}
	return updatedBdg, nil
}

func (s *BudgetServerImpl) DeleteBudget(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, UserID := ProtoBudgetReqToInts(req)
	deletedBdg, err := s.bdgUC.DeleteBudget(ctx, budgetID, UserID)
	if err != nil {
		return nil, err
	}
	return deletedBdg, nil
}

