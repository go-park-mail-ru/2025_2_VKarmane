package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Service struct {
	repo  FinanceRepository
	es    *elasticsearch.Client
	clock clock.Clock
}

func NewService(repo FinanceRepository, es *elasticsearch.Client, clck clock.Clock) *Service {
	return &Service{
		repo:  repo,
		es:    es,
		clock: clck,
	}
}

// Account methods
func (s *Service) GetAccountsByUser(ctx context.Context, userID int) (*finpb.ListAccountsResponse, error) {
	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	protoAccounts := make([]*finpb.Account, 0, len(accounts))
	for _, acc := range accounts {
		protoAccounts = append(protoAccounts, accountToProto(acc))
	}

	return &finpb.ListAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}

func (s *Service) GetAccountByID(ctx context.Context, userID, accountID int) (*finpb.Account, error) {
	account, err := s.repo.GetAccountByID(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	return accountToProto(account), nil
}

func (s *Service) CreateAccount(ctx context.Context, req finmodels.CreateAccountRequest) (*finpb.Account, error) {
	account := finmodels.Account{
		Balance:     req.Balance,
		Type:        req.Type,
		Name:        req.Name,
		CurrencyID:  req.CurrencyID,
		Description: &req.Description,
		CreatedAt:   s.clock.Now(),
		UpdatedAt:   s.clock.Now(),
	}

	createdAcc, err := s.repo.CreateAccount(ctx, account, req.UserID)
	if err != nil {
		return nil, err
	}

	return accountToProto(createdAcc), nil
}

func (s *Service) UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (*finpb.Account, error) {
	updatedAcc, err := s.repo.UpdateAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return accountToProto(updatedAcc), nil
}

func (s *Service) DeleteAccount(ctx context.Context, userID, accountID int) (*finpb.Account, error) {
	deletedAcc, err := s.repo.DeleteAccount(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	return accountToProto(deletedAcc), nil
}

func (s *Service) AddUserToAccount(ctx context.Context, userID, accountID int) (*finpb.SharingsResponse, error) {
	sh, err := s.repo.AddUserToAccount(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	return SharingToProto(sh), nil
}

// Operation methods
func (s *Service) GetOperationsByAccount(ctx context.Context, req []byte) (*finpb.ListOperationsResponse, error) {
	res, err := s.es.Search(
		s.es.Search.WithContext(context.Background()),
		s.es.Search.WithIndex("transactions"),
		s.es.Search.WithBody(bytes.NewReader(req)),
		s.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var esResp finmodels.ElasticsearchResponse
	if err := json.Unmarshal(body, &esResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ES: %w", err)
	}

	result := &finpb.ListOperationsResponse{}

	for _, hit := range esResp.Hits.Hits {
		op := convertToOperation(hit.Source)
		result.Operations = append(result.Operations, op)
	}

	return result, nil
}

func (s *Service) GetOperationByID(ctx context.Context, accID, opID int) (*finpb.Operation, error) {
	operation, err := s.repo.GetOperationByID(ctx, accID, opID)
	if err != nil {
		return nil, err
	}
	return operationToProto(operation), nil
}

func (s *Service) CreateOperation(ctx context.Context, req finmodels.CreateOperationRequest, accountID int) (*finpb.Operation, error) {
	var categoryID int
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	operationDate := s.clock.Now()
	if req.Date != nil {
		operationDate = *req.Date
	}

	op := finmodels.Operation{
		AccountID:   accountID,
		CategoryID:  categoryID,
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		Sum:         req.Sum,
		Status:      finmodels.OperationFinished,
		CreatedAt:   s.clock.Now(),
		Date:        operationDate,
	}

	createdOp, err := s.repo.CreateOperation(ctx, op)
	if err != nil {
		return nil, err
	}

	// Update account balance
	// account, err := s.repo.GetAccountByID(ctx, req.UserID, accountID)
	// if err != nil {
	// 	return nil, err
	// }

	// var newBalance float64
	// if req.Type == finmodels.OperationIncome {
	// 	newBalance = account.Balance + req.Sum
	// } else {
	// 	newBalance = account.Balance - req.Sum
	// }

	// if err := s.repo.UpdateAccountBalance(ctx, accountID, newBalance); err != nil {
	// 	return nil, err
	// }

	return operationToProto(createdOp), nil
}

func (s *Service) UpdateOperation(ctx context.Context, req finmodels.UpdateOperationRequest) (*finpb.Operation, error) {
	updatedOp, err := s.repo.UpdateOperation(ctx, req, req.AccountID, req.OperationID)
	if err != nil {
		return nil, err
	}
	return operationToProto(updatedOp), nil
}

func (s *Service) DeleteOperation(ctx context.Context, accID, opID int) (*finpb.Operation, error) {
	deletedOp, err := s.repo.DeleteOperation(ctx, accID, opID)
	if err != nil {
		return nil, err
	}
	return operationToProto(deletedOp), nil
}

// Category methods
func (s *Service) CreateCategory(ctx context.Context, req finmodels.CreateCategoryRequest) (*finpb.Category, error) {
	logoHashedID := req.LogoHashedID
	if logoHashedID == "" {
		logoHashedID = "c1dfd96eea8cc2b62785275bca38ac261256e278"
	}

	category := finmodels.Category{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		LogoHashedID: logoHashedID,
		CreatedAt:    s.clock.Now(),
		UpdatedAt:    s.clock.Now(),
	}

	createdCat, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return CategoryToProto(createdCat), nil
}

func (s *Service) GetCategoriesByUser(ctx context.Context, userID int) (*finpb.ListCategoriesResponse, error) {
	categories, err := s.repo.GetCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	protoCats := make([]*finpb.Category, 0, len(categories))
	for _, cat := range categories {
		protoCats = append(protoCats, CategoryToProto(cat))
	}

	return &finpb.ListCategoriesResponse{
		Categories: protoCats,
	}, nil
}

func (s *Service) GetCategoriesWithStatsByUser(ctx context.Context, userID int) (*finpb.ListCategoriesWithStatsResponse, error) {
	categories, err := s.repo.GetCategoriesWithStatsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	protoCats := make([]*finpb.CategoryWithStats, 0, len(categories))
	for _, cat := range categories {
		protoCats = append(protoCats, CategoryWithStatsToProto(cat.Category, cat.OperationsCount))
	}

	return &finpb.ListCategoriesWithStatsResponse{
		Categories: protoCats,
	}, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, userID, categoryID int) (*finpb.CategoryWithStats, error) {
	category, err := s.repo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return nil, err
	}
	stats, err := s.repo.GetCategoryStats(ctx, userID, categoryID)
	if err != nil {
		return nil, err
	}

	return CategoryWithStatsToProto(category, stats), nil
}

func (s *Service) UpdateCategory(ctx context.Context, category finmodels.Category) (*finpb.Category, error) {
	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	updatedCat, err := s.repo.GetCategoryByID(ctx, category.UserID, category.ID)
	if err != nil {
		return nil, err
	}

	return CategoryToProto(updatedCat), nil
}

func (s *Service) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	return s.repo.DeleteCategory(ctx, userID, categoryID)
}
