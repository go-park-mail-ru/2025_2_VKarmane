package impl

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	finerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/proto"
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
		if errors.Is(err, finerrors.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return account, nil
}

func (s *FinanceServerImpl) GetAccount(ctx context.Context, req *finpb.AccountRequest) (*finpb.Account, error) {
	account, err := s.financeUC.GetAccountByID(ctx, int(req.UserId), int(req.AccountId))
	if err != nil {
		if errors.Is(err, finerrors.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return account, nil
}

func (s *FinanceServerImpl) GetAccountsByUser(ctx context.Context, req *finpb.UserID) (*finpb.ListAccountsResponse, error) {
	accounts, err := s.financeUC.GetAccountsByUser(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return accounts, nil
}

func (s *FinanceServerImpl) UpdateAccount(ctx context.Context, req *finpb.UpdateAccountRequest) (*finpb.Account, error) {
	updateReq := protoToUpdateAccountRequest(req)
	account, err := s.financeUC.UpdateAccount(ctx, updateReq)
	if err != nil {
		if errors.Is(err, finerrors.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return account, nil
}

func (s *FinanceServerImpl) DeleteAccount(ctx context.Context, req *finpb.AccountRequest) (*finpb.Account, error) {
	account, err := s.financeUC.DeleteAccount(ctx, int(req.UserId), int(req.AccountId))
	if err != nil {
		if errors.Is(err, finerrors.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return account, nil
}

// Operation methods
func (s *FinanceServerImpl) CreateOperation(ctx context.Context, req *finpb.CreateOperationRequest) (*finpb.Operation, error) {
	createReq := protoToCreateOperationRequest(req)
	operation, err := s.financeUC.CreateOperation(ctx, createReq, int(req.AccountId))
	if err != nil {
		if errors.Is(err, finerrors.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return operation, nil
}

func (s *FinanceServerImpl) GetOperation(ctx context.Context, req *finpb.OperationRequest) (*finpb.Operation, error) {
	operation, err := s.financeUC.GetOperationByID(ctx, int(req.AccountId), int(req.OperationId))
	if err != nil {
		if errors.Is(err, finerrors.ErrOperationNotFound) {
			return nil, status.Error(codes.NotFound, "operation not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return operation, nil
}

func (s *FinanceServerImpl) GetOperationsByAccount(ctx context.Context, req *finpb.AccountRequest) (*finpb.ListOperationsResponse, error) {
	operations, err := s.financeUC.GetOperationsByAccount(ctx, int(req.AccountId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return operations, nil
}

func (s *FinanceServerImpl) UpdateOperation(ctx context.Context, req *finpb.UpdateOperationRequest) (*finpb.Operation, error) {
	updateReq := protoToUpdateOperationRequest(req)
	operation, err := s.financeUC.UpdateOperation(ctx, updateReq)
	if err != nil {
		if errors.Is(err, finerrors.ErrOperationNotFound) {
			return nil, status.Error(codes.NotFound, "operation not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return operation, nil
}

func (s *FinanceServerImpl) DeleteOperation(ctx context.Context, req *finpb.OperationRequest) (*finpb.Operation, error) {
	operation, err := s.financeUC.DeleteOperation(ctx, int(req.AccountId), int(req.OperationId))
	if err != nil {
		if errors.Is(err, finerrors.ErrOperationNotFound) {
			return nil, status.Error(codes.NotFound, "operation not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return operation, nil
}

// Category methods
func (s *FinanceServerImpl) CreateCategory(ctx context.Context, req *finpb.CreateCategoryRequest) (*finpb.Category, error) {
	createReq := protoToCreateCategoryRequest(req)
	category, err := s.financeUC.CreateCategory(ctx, createReq)
	if err != nil {
		if errors.Is(err, finerrors.ErrCategoryExists) {
			return nil, status.Error(codes.AlreadyExists, "category already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return category, nil
}

func (s *FinanceServerImpl) GetCategory(ctx context.Context, req *finpb.CategoryRequest) (*finpb.Category, error) {
	category, err := s.financeUC.GetCategoryByID(ctx, int(req.UserId), int(req.CategoryId))
	if err != nil {
		if errors.Is(err, finerrors.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "category not found")
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
		ID:           int(updateReq.CategoryId),
		UserID:       int(updateReq.UserId),
		Name:         getStringValue(updateReq.Name),
		Description:  getStringValue(updateReq.Description),
		LogoHashedID: getStringValue(updateReq.LogoHashedId),
	}
	updatedCategory, err := s.financeUC.UpdateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, finerrors.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return updatedCategory, nil
}

func (s *FinanceServerImpl) DeleteCategory(ctx context.Context, req *finpb.CategoryRequest) (*finpb.Category, error) {
	err := s.financeUC.DeleteCategory(ctx, int(req.UserId), int(req.CategoryId))
	if err != nil {
		if errors.Is(err, finerrors.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		if errors.Is(err, finerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	// Return empty category on successful deletion
	return &finpb.Category{}, nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

