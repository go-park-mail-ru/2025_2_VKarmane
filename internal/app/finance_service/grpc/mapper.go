package grpc

import (
	"time"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
)

func protoToCreateAccountRequest(req *finpb.CreateAccountRequest) finmodels.CreateAccountRequest {
	var description string
	if req.Description != nil {
		description = *req.Description
	}

	return finmodels.CreateAccountRequest{
		UserID:      int(req.UserId),
		Balance:     req.Balance,
		Type:        req.Type,
		CurrencyID:  int(req.CurrencyId),
		Name:        req.Name,
		Description: description,
	}
}

func protoToUpdateAccountRequest(req *finpb.UpdateAccountRequest) finmodels.UpdateAccountRequest {
	return finmodels.UpdateAccountRequest{
		UserID:      int(req.UserId),
		AccountID:   int(req.AccountId),
		Balance:     req.Balance,
		Name:        req.Name,
		Description: req.Description,
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
		UserID:      int(req.UserId),
		AccountID:   int(req.AccountId),
		CategoryID:  categoryID,
		Type:        finmodels.OperationType(req.Type),
		Name:        req.Name,
		Description: req.Description,
		Sum:         req.Sum,
		Date:        date,
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

func protoToCategoryRequest(req *finpb.CategoryReportRequest) finmodels.CategoryReportRequest {
	return finmodels.CategoryReportRequest{
		UserID: int(req.UserId),
		Start:  req.Start.AsTime(),
		End:    req.End.AsTime(),
	}
}
