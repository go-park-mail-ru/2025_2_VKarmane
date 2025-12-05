package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	finerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type FinanceServerImpl struct {
	finpb.UnimplementedFinanceServiceServer
	financeUC FinanceUseCase
}

func NewFinanceServer(financeUC FinanceUseCase) *FinanceServerImpl {
	return &FinanceServerImpl{
		financeUC: financeUC,
	}
}

// Account methods
func (s *FinanceServerImpl) CreateAccount(ctx context.Context, req *finpb.CreateAccountRequest) (*finpb.Account, error) {
	createReq := protoToCreateAccountRequest(req)
	account, err := s.financeUC.CreateAccount(ctx, createReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to create account", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to create account, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return account, nil
}

func (s *FinanceServerImpl) GetAccount(ctx context.Context, req *finpb.AccountRequest) (*finpb.Account, error) {
	account, err := s.financeUC.GetAccountByID(ctx, int(req.UserId), int(req.AccountId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get account", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get account, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return account, nil
}

func (s *FinanceServerImpl) GetAccountsByUser(ctx context.Context, req *finpb.UserID) (*finpb.ListAccountsResponse, error) {
	accounts, err := s.financeUC.GetAccountsByUser(ctx, int(req.UserId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get accounts", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get accounts, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return accounts, nil
}

func (s *FinanceServerImpl) UpdateAccount(ctx context.Context, req *finpb.UpdateAccountRequest) (*finpb.Account, error) {
	updateReq := protoToUpdateAccountRequest(req)
	account, err := s.financeUC.UpdateAccount(ctx, updateReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to update account", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to update account, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return account, nil
}

func (s *FinanceServerImpl) DeleteAccount(ctx context.Context, req *finpb.AccountRequest) (*finpb.Account, error) {
	account, err := s.financeUC.DeleteAccount(ctx, int(req.UserId), int(req.AccountId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to delete account", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to delete account, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return account, nil
}

// Operation methods
func (s *FinanceServerImpl) CreateOperation(ctx context.Context, req *finpb.CreateOperationRequest) (*finpb.Operation, error) {
	createReq := protoToCreateOperationRequest(req)
	operation, err := s.financeUC.CreateOperation(ctx, createReq, int(req.AccountId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to create operation", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to create operation, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return operation, nil
}

func (s *FinanceServerImpl) AddUserToAccounnt(ctx context.Context, req *finpb.AccountRequest) (*finpb.SharingsResponse, error) {
	sharing, err := s.financeUC.AddUserToAccount(ctx, int(req.UserId), int(req.AccountId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to add user to account", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to add user to account, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeAccountNotFound))
	}
	return sharing, nil
}

func (s *FinanceServerImpl) GetOperation(ctx context.Context, req *finpb.OperationRequest) (*finpb.Operation, error) {
	operation, err := s.financeUC.GetOperationByID(ctx, int(req.AccountId), int(req.OperationId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get operation", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get operation, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return operation, nil
}

func (s *FinanceServerImpl) GetOperationsByAccount(ctx context.Context, req *finpb.OperationsByAccountAndFiltersRequest) (*finpb.ListOperationsResponse, error) {
	logger := logger.FromContext(ctx)

	var date string
	if req.Date != nil {
		date = req.Date.AsTime().Format(time.RFC3339Nano)
	} else {
		date = ""
	}

	operations, err := s.financeUC.GetOperationsByAccount(ctx, int(req.AccountId), int(req.CategoryId), req.Name, req.OperationType, req.AccountType, date)
	if err != nil {
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get operations", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get operations, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return operations, nil
}

func (s *FinanceServerImpl) UpdateOperation(ctx context.Context, req *finpb.UpdateOperationRequest) (*finpb.Operation, error) {
	updateReq := protoToUpdateOperationRequest(req)
	operation, err := s.financeUC.UpdateOperation(ctx, updateReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to update operation", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to update operation, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return operation, nil
}

func (s *FinanceServerImpl) DeleteOperation(ctx context.Context, req *finpb.OperationRequest) (*finpb.Operation, error) {
	operation, err := s.financeUC.DeleteOperation(ctx, int(req.AccountId), int(req.OperationId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to delete operation", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to delete operation, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return operation, nil
}

// Category methods
func (s *FinanceServerImpl) CreateCategory(ctx context.Context, req *finpb.CreateCategoryRequest) (*finpb.Category, error) {
	createReq := protoToCreateCategoryRequest(req)
	category, err := s.financeUC.CreateCategory(ctx, createReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to create category", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to create category, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return category, nil
}

func (s *FinanceServerImpl) GetCategory(ctx context.Context, req *finpb.CategoryRequest) (*finpb.CategoryWithStats, error) {
	category, err := s.financeUC.GetCategoryByID(ctx, int(req.UserId), int(req.CategoryId))
	if err != nil {
		logger := logger.FromContext(ctx)
		if errors.Is(err, finerrors.ErrCategoryNotFound) {
			if logger != nil {
				logger.Error("Failed to get category", "error", err)
			}
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if logger != nil {
			logger.Error("Failed to get category, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return category, nil
}

func (s *FinanceServerImpl) GetCategoriesByUser(ctx context.Context, req *finpb.UserID) (*finpb.ListCategoriesResponse, error) {
	categories, err := s.financeUC.GetCategoriesByUser(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return categories, nil
}

func (s *FinanceServerImpl) GetCategoryByName(ctx context.Context, req *finpb.CategoryByNameRequest) (*finpb.CategoryWithStats, error) {
	category, err := s.financeUC.GetCategoryByName(ctx, int(req.UserId), req.CategoryName)
	if err != nil {
		logger := logger.FromContext(ctx)
		if errors.Is(err, finerrors.ErrCategoryNotFound) {
			if logger != nil {
				logger.Error("Failed to get category", "error", err)
			}
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if logger != nil {
			logger.Error("Failed to get category, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return category, nil
}

func (s *FinanceServerImpl) GetCategoriesWithStatsByUser(ctx context.Context, req *finpb.UserID) (*finpb.ListCategoriesWithStatsResponse, error) {
	categories, err := s.financeUC.GetCategoriesWithStatsByUser(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return categories, nil
}

func (s *FinanceServerImpl) UpdateCategory(ctx context.Context, req *finpb.UpdateCategoryRequest) (*finpb.Category, error) {
	updateReq := protoToUpdateCategoryRequest(req)
	category := finmodels.Category{
		ID:           int(updateReq.CategoryID),
		UserID:       int(updateReq.UserID),
		Name:         getStringValue(updateReq.Name),
		Description:  getStringValue(updateReq.Description),
		LogoHashedID: getStringValue(updateReq.LogoHashedID),
	}
	updatedCategory, err := s.financeUC.UpdateCategory(ctx, category)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to update category", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to update category, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}

	return updatedCategory, nil
}

func (s *FinanceServerImpl) DeleteCategory(ctx context.Context, req *finpb.CategoryRequest) (*finpb.Category, error) {
	err := s.financeUC.DeleteCategory(ctx, int(req.UserId), int(req.CategoryId))
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range finerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to delete category", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to delete category, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return &finpb.Category{}, nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
