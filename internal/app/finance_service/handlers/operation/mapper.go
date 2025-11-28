package operation

import (
	"net/url"
	"strconv"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AccountAndUserIDToProtoID(accountID, userID int) *finpb.AccountRequest {
	return &finpb.AccountRequest{
		UserId:    int32(userID),
		AccountId: int32(accountID),
	}
}

func MapOperationInListToResponse(op *finpb.OperationInList) models.OperationInListResponse {
	var createdAt time.Time
	if op.CreatedAt != nil {
		createdAt = op.CreatedAt.AsTime()
	}

	var date time.Time
	if op.Date != nil {
		date = op.Date.AsTime()
	}

	return models.OperationInListResponse{
		ID:               int(op.Id),
		AccountID:        int(op.AccountId),
		CategoryID:       int(op.CategoryId),
		CategoryName:     op.CategoryName,
		Name:             op.Name,
		Type:             op.Type,
		Description:      op.Description,
		CategoryLogo:     op.CategoryLogo,
		CategoryHashedID: op.CategoryLogoHashedId,
		Sum:              op.Sum,
		CurrencyID:       int(op.CurrencyId),
		CreatedAt:        createdAt,
		Date:             date,
	}
}

func MapListOperationsResponse(resp *finpb.ListOperationsResponse) []models.OperationInListResponse {
	if resp == nil {
		return nil
	}

	ops := make([]models.OperationInListResponse, 0, len(resp.Operations))
	for _, op := range resp.Operations {
		ops = append(ops, MapOperationInListToResponse(op))
	}
	return ops
}

func CreateOperationRequestToProto(req models.CreateOperationRequest, userID, accID int) *finpb.CreateOperationRequest {
	var categoryID int32
	if req.CategoryID != nil {
		categoryID = int32(*req.CategoryID)
	}

	var date *timestamppb.Timestamp
	if req.Date != nil {
		date = timestamppb.New(*req.Date)
	}

	return &finpb.CreateOperationRequest{
		UserId:      int32(userID),
		AccountId:   int32(accID),
		CategoryId:  &categoryID,
		Type:        string(req.Type),
		Description: req.Description,
		Name:        req.Name,
		Sum:         req.Sum,
		Date:        date,
	}
}

func ProtoOperationToResponse(op *finpb.Operation) models.OperationResponse {
	return models.OperationResponse{
		ID:           int(op.Id),
		AccountID:    int(op.AccountId),
		CategoryID:   int(op.CategoryId),
		CategoryName: op.CategoryName,
		Type:         string(op.Type),
		Status:       string(op.Status),
		Description:  op.Description,
		ReceiptURL:   op.ReceiptUrl,
		Name:         op.Name,
		Sum:          op.Sum,
		CurrencyID:   int(op.CurrencyId),
		CreatedAt:    op.CreatedAt.AsTime(),
		Date:         op.Date.AsTime(),
	}
}

func OperationAndUserIDToProtoID(opID, accID, userID int) *finpb.OperationRequest {
	return &finpb.OperationRequest{
		AccountId:   int32(accID),
		UserId:      int32(userID),
		OperationId: int32(opID),
	}
}

func UpdateOperationRequestToProto(req models.UpdateOperationRequest, userID, accID, opID int) *finpb.UpdateOperationRequest {
	var categoryID int32
	if req.CategoryID != nil {
		categoryID = int32(*req.CategoryID)
	}

	var date *timestamppb.Timestamp
	if req.CreatedAt != nil {
		date = timestamppb.New(*req.CreatedAt)
	}

	return &finpb.UpdateOperationRequest{
		UserId:      int32(userID),
		AccountId:   int32(accID),
		OperationId: int32(opID),
		CategoryId:  &categoryID,
		Description: req.Description,
		Name:        req.Name,
		Sum:         req.Sum,
		CreatedAt:   date,
	}
}

func MapSearchOperationsRequest(q url.Values, userID, accID int) *finpb.SearchOperationsRequest {
	req := &finpb.SearchOperationsRequest{}

	req.AccountId = int32(accID)
	req.UserId = int32(userID)

	if s := q.Get("category_id"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			req.CategoryId = int32(v)
		}
	}

	if s := q.Get("operation_type"); s != "" {
		req.OperationType = s
	}

	if s := q.Get("account_type"); s != "" {
		req.AccountType = s
	}

	if s := q.Get("title"); s != "" {
		req.Name = s
	}
	if s := q.Get("date_from"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			req.CreatedAt = timestamppb.New(t)
		}
	}

	return req
}
