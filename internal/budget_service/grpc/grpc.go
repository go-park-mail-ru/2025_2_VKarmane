package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/errors"
	budgetpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetServiceServer struct {
	bdgUC BudgetUseCase
	budgetpb.UnimplementedBudgetServiceServer
}

func NewBudgetServer(bdgUC BudgetUseCase) *BudgetServiceServer {
	return &BudgetServiceServer{bdgUC: bdgUC}
}

func (s *BudgetServiceServer) CreateBudget(ctx context.Context, req *budgetpb.CreateBudgetRequest) (*budgetpb.Budget, error) {
	budget, userID := ProtoCreateRequestToModel(req)
	createdBdg, err := s.bdgUC.CreateBudget(ctx, budget, userID)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range bdgerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to create budget", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to create budget, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return createdBdg, nil
}

func (s *BudgetServiceServer) GetBudgetByID(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, userID := ProtoBudgetReqToInts(req)
	bdg, err := s.bdgUC.GetBudget(ctx, budgetID, userID)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range bdgerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get budget", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get budget, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return bdg, nil
}

func (s *BudgetServiceServer) GetListBudgets(ctx context.Context, req *budgetpb.UserID) (*budgetpb.ListBudgetsResponse, error) {
	userID := ProtoIDToInt(req)
	bdgList, err := s.bdgUC.GetBudgets(ctx, userID)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range bdgerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get budgets", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get budgets, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return bdgList, nil
}

func (s *BudgetServiceServer) UpdateBudget(ctx context.Context, req *budgetpb.UpdateBudgetRequest) (*budgetpb.Budget, error) {
	budget := ProtoUpdateRequestToModel(req)
	updatedBdg, err := s.bdgUC.UpdateBudget(ctx, budget)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range bdgerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to update budget", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to update budget, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return updatedBdg, nil
}

func (s *BudgetServiceServer) DeleteBudget(ctx context.Context, req *budgetpb.BudgetRequest) (*budgetpb.Budget, error) {
	budgetID, userID := ProtoBudgetReqToInts(req)
	deletedBdg, err := s.bdgUC.DeleteBudget(ctx, budgetID, userID)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range bdgerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to delete budget", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to delete budget, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return deletedBdg, nil
}
