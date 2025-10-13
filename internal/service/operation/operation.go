package operation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

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
    accs := s.accountRepo.GetAccountsByUser(ctx, userID)
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
		return []models.Operation{}, fmt.Errorf("forbidden: account does not belong to user")
	}
	ops := s.operationRepo.GetOperationsByAccount(ctx, accID)

	return ops, nil
}

func (s *Service) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, fmt.Errorf("forbidden: account does not belong to user")
	}
	ops := s.operationRepo.GetOperationsByAccount(ctx, accID)

	for _, op := range ops {
		if op.ID == opID {
			return op, nil
		}
	}

	return models.Operation{}, nil
}


func (s *Service) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, fmt.Errorf("forbidden: account does not belong to user")
	}
	op := models.Operation{
		ID: 0,
		AccountID: req.AccountID,
		CategoryID: req.CategoryID,
		Type: req.Type,
		Status: models.OperationFinished,
		Description: req.Description,
		ReceiptURL: "11111111111",
		Name: req.Name,
		Sum: req.Sum,
		CurrencyID: 1,
		CreatedAt: time.Now(),
	}

	createdOp := s.operationRepo.CreateOperation(ctx, op)
	
	return createdOp, nil
}



func (s *Service) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, fmt.Errorf("forbidden: account does not belong to user")
	}
	updatedOp := s.operationRepo.UpdateOperation(ctx, req, accID, opID)
	
	return updatedOp, nil
}


func (s *Service) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	if !(s.CheckAccountOwnership(ctx, accID)) {
		return models.Operation{}, fmt.Errorf("forbidden: account does not belong to user")
	}
	deletedOp := s.operationRepo.DeleteOperation(ctx, accID, opID)
	
	return deletedOp, nil
}