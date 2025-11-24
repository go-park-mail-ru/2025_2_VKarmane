package service

import (
	"time"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func accountToProto(account finmodels.Account) *finpb.Account {
	return &finpb.Account{
		Id:         int32(account.ID),
		Balance:    account.Balance,
		Type:       account.Type,
		CurrencyId: int32(account.CurrencyID),
		CreatedAt:  timestamppb.New(account.CreatedAt),
		UpdatedAt:  timestamppb.New(account.UpdatedAt),
	}
}

func operationToProto(operation finmodels.Operation) *finpb.Operation {
	return &finpb.Operation{
		Id:           int32(operation.ID),
		AccountId:    int32(operation.AccountID),
		CategoryId:   int32(operation.CategoryID),
		CategoryName: operation.CategoryName,
		Type:         string(operation.Type),
		Status:       string(operation.Status),
		Description:  operation.Description,
		ReceiptUrl:   operation.ReceiptURL,
		Name:         operation.Name,
		Sum:          operation.Sum,
		CurrencyId:   int32(operation.CurrencyID),
		CreatedAt:    timestamppb.New(operation.CreatedAt),
		Date:         timestamppb.New(operation.Date),
	}
}

func operationInListToProto(operation finmodels.OperationInList) *finpb.OperationInList {
	return &finpb.OperationInList{
		Id:                 int32(operation.ID),
		AccountId:          int32(operation.AccountID),
		CategoryId:         int32(operation.CategoryID),
		CategoryName:       operation.CategoryName,
		Type:               string(operation.Type),
		Description:        operation.Description,
		Name:               operation.Name,
		CategoryLogoHashedId: operation.CategoryLogoHashedID,
		CategoryLogo:       operation.CategoryLogo,
		Sum:                operation.Sum,
		CurrencyId:         int32(operation.CurrencyID),
		CreatedAt:          timestamppb.New(operation.CreatedAt),
		Date:               timestamppb.New(operation.Date),
	}
}

func categoryToProto(category finmodels.Category) *finpb.Category {
	return &finpb.Category{
		Id:           int32(category.ID),
		UserId:       int32(category.UserID),
		Name:         category.Name,
		Description:  category.Description,
		LogoHashedId: category.LogoHashedID,
		LogoUrl:      category.LogoURL,
		CreatedAt:    timestamppb.New(category.CreatedAt),
		UpdatedAt:    timestamppb.New(category.UpdatedAt),
	}
}

func categoryWithStatsToProto(category finmodels.CategoryWithStats) *finpb.CategoryWithStats {
	return &finpb.CategoryWithStats{
		Category:        categoryToProto(category.Category),
		OperationsCount: int32(category.OperationsCount),
	}
}

func protoToCreateAccountRequest(req *finpb.CreateAccountRequest) finmodels.CreateAccountRequest {
	return finmodels.CreateAccountRequest{
		UserID:     int(req.UserId),
		Balance:    req.Balance,
		Type:       req.Type,
		CurrencyID: int(req.CurrencyId),
	}
}

func protoToUpdateAccountRequest(req *finpb.UpdateAccountRequest) finmodels.UpdateAccountRequest {
	return finmodels.UpdateAccountRequest{
		UserID:    int(req.UserId),
		AccountID: int(req.AccountId),
		Balance:   req.Balance,
	}
}

func protoToCreateOperationRequest(req *finpb.CreateOperationRequest) finmodels.CreateOperationRequest {
	var categoryID *int
	if req.CategoryId != nil {
		id := int(*req.CategoryId)
		categoryID = &id
	}

	var date *time.Time
	if req.Date != nil {
		t := req.Date.AsTime()
		date = &t
	}

	return finmodels.CreateOperationRequest{
		UserID:     int(req.UserId),
		AccountID:  int(req.AccountId),
		CategoryID: categoryID,
		Type:       finmodels.OperationType(req.Type),
		Name:       req.Name,
		Description: req.Description,
		Sum:        req.Sum,
		Date:       date,
	}
}

func protoToUpdateOperationRequest(req *finpb.UpdateOperationRequest) finmodels.UpdateOperationRequest {
	var categoryID *int
	if req.CategoryId != nil {
		id := int(*req.CategoryId)
		categoryID = &id
	}

	var name *string
	if req.Name != nil {
		name = req.Name
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	var sum *float64
	if req.Sum != nil {
		sum = req.Sum
	}

	var createdAt *time.Time
	if req.CreatedAt != nil {
		t := req.CreatedAt.AsTime()
		createdAt = &t
	}

	return finmodels.UpdateOperationRequest{
		UserID:      int(req.UserId),
		AccountID:   int(req.AccountId),
		OperationID: int(req.OperationId),
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Sum:         sum,
		CreatedAt:   createdAt,
	}
}

func protoToCreateCategoryRequest(req *finpb.CreateCategoryRequest) finmodels.CreateCategoryRequest {
	return finmodels.CreateCategoryRequest{
		UserID:       int(req.UserId),
		Name:         req.Name,
		Description:  req.Description,
		LogoHashedID: req.LogoHashedId,
	}
}

func protoToUpdateCategoryRequest(req *finpb.UpdateCategoryRequest) finmodels.UpdateCategoryRequest {
	var name *string
	if req.Name != nil {
		name = req.Name
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	var logoHashedID *string
	if req.LogoHashedId != nil {
		logoHashedID = req.LogoHashedId
	}

	return finmodels.UpdateCategoryRequest{
		UserID:       int(req.UserId),
		CategoryID:   int(req.CategoryId),
		Name:         name,
		Description:  description,
		LogoHashedID: logoHashedID,
	}
}

