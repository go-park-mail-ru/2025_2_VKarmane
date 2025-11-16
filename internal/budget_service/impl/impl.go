package budget

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	budgetpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/errors"
)

type BudgetServerImpl struct {
	bdgUC BudgetUseCase
	budgetpb.UnimplementedBudgetServiceServer
}

func NewBudgetServer(bdgUC BudgetUseCase) *BudgetServerImpl {
	return &BudgetServerImpl{bdgUC: bdgUC}
}

func (s *BudgetServerImpl) CreateBudget(ctx context.Context, req *budgetpb.CreateBudgetRequest) (*budgetpb.Budget, error) {
	budget, userID := ProtoCreateRequestToModel(req)
	createdBdg, err := s.bdgUC.CreateBudget(ctx, budget, userID)
	if err != nil {
		if errors.Is(err, bdgerrors.ErrBudgetExists) {
			return nil, status.Error(codes.AlreadyExists, "budget already exists")
		}
		if errors.Is(err, bdgerrors.ErrInavlidData) {
			return nil, status.Error(codes.InvalidArgument, "invalid budget data")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return createdBdg, nil
}

func (s *BudgetServerImpl) GetBudgetByID(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, userID := ProtoBudgetReqToInts(req)
	bdg, err := s.bdgUC.GetBudget(ctx, budgetID, userID)
	if err != nil {
		if errors.Is(err, bdgerrors.ErrBudgetNotFound) {
			return nil, status.Error(codes.NotFound, "budget not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return bdg, nil
}

func (s *BudgetServerImpl) GetListBudgets(ctx context.Context, req *budgetpb.UserID) (*budgetpb.ListBudgetsResponse, error) {
	userID := ProtoIDToInt(req)
	bdgList, err := s.bdgUC.GetBudgets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return bdgList, nil
}

func (s *BudgetServerImpl) UpdateBudget(ctx context.Context, req *budgetpb.UpdateBudgetRequest) (*budgetpb.Budget, error) {
	budget := ProtoUpdateRequestToModel(req)
	updatedBdg, err := s.bdgUC.UpdateBudget(ctx, budget)
	if err != nil {
		if errors.Is(err, bdgerrors.ErrBudgetNotFound) {
			return nil, status.Error(codes.NotFound, "budget not found")
		}
		if errors.Is(err, bdgerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return updatedBdg, nil
}

func (s *BudgetServerImpl) DeleteBudget(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, userID := ProtoBudgetReqToInts(req)
	deletedBdg, err := s.bdgUC.DeleteBudget(ctx, budgetID, userID)
	if err != nil {
		if errors.Is(err, bdgerrors.ErrBudgetNotFound) {
			return nil, status.Error(codes.NotFound, "budget not found")
		}
		if errors.Is(err, bdgerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return deletedBdg, nil
}
