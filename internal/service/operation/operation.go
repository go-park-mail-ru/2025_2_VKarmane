package operation

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

var ErrForbidden = errors.New("forbidden")

type Service struct {
	accountRepo   AccountRepository
	operationRepo OperationRepository
	clock         clock.Clock
}

func NewService(accountRepo AccountRepository, operationRepo OperationRepository, clck clock.Clock) *Service {
	return &Service{
		accountRepo:   accountRepo,
		operationRepo: operationRepo,
		clock:         clck,
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
	if !s.CheckAccountOwnership(ctx, accID) {
		return []models.Operation{}, ErrForbidden
	}
	ops, err := s.operationRepo.GetOperationsByAccount(ctx, accID)
	if err != nil {
		return []models.Operation{}, pkgerrors.Wrap(err, "Failed to get account operations")
	}
	return ops, nil
}

func (s *Service) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Operation{}, ErrForbidden
	}
	op, err := s.operationRepo.GetOperationByID(ctx, accID, opID)
	if err != nil {
		return models.Operation{}, pkgErrors.Wrap(err, "Failed to get operation by ID")
	}
	return op, nil
}

func (s *Service) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Operation{}, ErrForbidden
	}

	var categoryID int
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	operationDate := time.Now()
	if req.Date != nil {
		operationDate = *req.Date
	}

	operation := models.Operation{
		AccountID:   accID,
		CategoryID:  categoryID,
		Type:        req.Type,
		Status:      models.OperationFinished,
		Description: req.Description,
		ReceiptURL:  "",
		Name:        req.Name,
		Sum:         req.Sum,
		CurrencyID:  1, // TODO: получать из аккаунта или запроса
		CreatedAt:   time.Now(),
		Date:        operationDate,
	}

	createdOp, err := s.operationRepo.CreateOperation(ctx, operation)
	if err != nil {
		return models.Operation{}, pkgerrors.Wrap(err, "Failed to create operation")
	}

	return createdOp, nil
}

func (s *Service) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Operation{}, ErrForbidden
	}

	updatedOp, err := s.operationRepo.UpdateOperation(ctx, req, accID, opID)
	if err != nil {
		return models.Operation{}, pkgerrors.Wrap(err, "Failed to update operation")
	}

	return updatedOp, nil
}

func (s *Service) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Operation{}, ErrForbidden
	}

	deletedOp, err := s.operationRepo.DeleteOperation(ctx, accID, opID)
	if err != nil {
		return models.Operation{}, pkgerrors.Wrap(err, "Failed to delete operation")
	}

	return deletedOp, nil
}
