package service

import (
	"time"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
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
		AccountType:  string(operation.AccountType),
		CurrencyId:   int32(operation.CurrencyID),
		CreatedAt:    timestamppb.New(operation.CreatedAt),
		Date:         timestamppb.New(operation.Date),
	}
}

func CategoryToProto(category finmodels.Category) *finpb.Category {
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

func CategoryWithStatsToProto(category finmodels.Category, operationsCount int) *finpb.CategoryWithStats {
	return &finpb.CategoryWithStats{
		Category:        CategoryToProto(category),
		OperationsCount: int32(operationsCount),
	}
}

func convertToOperation(src finmodels.ESHitSource) *finpb.OperationInList {
	return &finpb.OperationInList{
		Id:                   src.ID,
		AccountId:            src.AccountID,
		CategoryId:           src.CategoryID,
		CategoryName:         src.CategoryName,
		Type:                 src.Type,
		Description:          src.Description,
		Name:                 src.Name,
		CategoryLogoHashedId: src.CategoryLogoHashedId,
		CategoryLogo:         src.CategoryLogo,
		Sum:                  src.Sum,
		CurrencyId:           src.CurrencyId,
		AccountType:          src.AccountType,
		CreatedAt:            parseTime(src.CreatedAt),
		Date:                 parseTime(src.Date),
	}
}

func parseTime(s string) *timestamppb.Timestamp {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return nil
	}
	return timestamppb.New(t)
}
