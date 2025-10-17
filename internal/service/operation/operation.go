package operation

import (
	"context"
	"errors"
	"time"

	pkgErrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

var ErrForbidden = errors.New("forbidden")

type Service struct {
	accountRepo   AccountRepository
	operationRepo OperationRepository
}

func NewService(accountRepo AccountRepository, operationRepo OperationRepository) *Service {
	return &Service{
		accountRepo:   accountRepo,
		operationRepo: operationRepo,
	}
}

func (s *Service) CheckAccountOwnership(ctx context.Context, accID int) bool {
	userID, _ := middleware.GetUserIDFromContext(ctx)
	accs, err := s.accountRepo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return false
	}
	if len(accs) > 0 {
		for _, acc := range accs {
			if acc.ID == accID {
				return true
			}
		}
	}
	return false
}

func (s *Service) GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return []models.Operation{}, ErrForbidden
	}
	ops, err := s.operationRepo.GetOperationsByAccount(ctx, accID)
	if err != nil {
		return []models.Operation{}, pkgErrors.Wrap(err, "Failed to get account operations")
	}

	return ops, nil
}

func (s *Service) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, ErrForbidden
	}
	ops, err := s.operationRepo.GetOperationsByAccount(ctx, accID)
	if err != nil {
		return models.Operation{}, pkgErrors.Wrap(err, "Failed to get operation by id")
	}

	for _, op := range ops {
		if op.ID == opID {
			return op, nil
		}
	}

	return models.Operation{}, nil
}

func (s *Service) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, ErrForbidden
	}
	op := models.Operation{
		ID:          0,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Status:      models.OperationFinished,
		Description: req.Description,
		ReceiptURL:  "11111111111",
		Name:        req.Name,
		Sum:         req.Sum,
		CurrencyID:  1,
		CreatedAt:   time.Now(),
	}

	createdOp, err := s.operationRepo.CreateOperation(ctx, op)
	if err != nil {
		return models.Operation{}, pkgErrors.Wrap(err, "Failed to create operation")
	}

	return createdOp, nil
}

func (s *Service) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, ErrForbidden
	}
	updatedOp, err := s.operationRepo.UpdateOperation(ctx, req, accID, opID)
	if err != nil {
		return models.Operation{}, pkgErrors.Wrap(err, "Failed to update operation")
	}

	return updatedOp, nil
}

func (s *Service) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, ErrForbidden
	}
	deletedOp, err := s.operationRepo.DeleteOperation(ctx, accID, opID)
	if err != nil {
		return models.Operation{}, pkgErrors.Wrap(err, "Failed to delete operation")
	}

	return deletedOp, nil
}
