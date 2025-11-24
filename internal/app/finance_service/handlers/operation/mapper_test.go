package operation

import (
	"testing"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAccountAndUserIDToProtoID(t *testing.T) {
	res := AccountAndUserIDToProtoID(7, 10)
	require.Equal(t, int32(10), res.UserId)
	require.Equal(t, int32(7), res.AccountId)
}

func TestMapOperationInListToResponse(t *testing.T) {
	created := timestamppb.New(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
	date := timestamppb.New(time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC))

	op := &finpb.OperationInList{
		Id:                   1,
		AccountId:            2,
		CategoryId:           3,
		CategoryName:         "Food",
		Name:                 "Buy milk",
		Type:                 "expense",
		Description:          "Grocery",
		CategoryLogo:         "url",
		CategoryLogoHashedId: "hash",
		Sum:                  123.45,
		CurrencyId:           840,
		CreatedAt:            created,
		Date:                 date,
	}

	res := MapOperationInListToResponse(op)

	require.Equal(t, 1, res.ID)
	require.Equal(t, 2, res.AccountID)
	require.Equal(t, 3, res.CategoryID)
	require.Equal(t, "Food", res.CategoryName)
	require.Equal(t, "Buy milk", res.Name)
	require.Equal(t, "expense", res.Type)
	require.Equal(t, "Grocery", res.Description)
	require.Equal(t, "url", res.CategoryLogo)
	require.Equal(t, "hash", res.CategoryHashedID)
	require.Equal(t, 123.45, res.Sum)
	require.Equal(t, 840, res.CurrencyID)
	require.Equal(t, created.AsTime(), res.CreatedAt)
	require.Equal(t, date.AsTime(), res.Date)
}

func TestMapListOperationsResponse(t *testing.T) {
	op1 := &finpb.OperationInList{Id: 1, Name: "A"}
	op2 := &finpb.OperationInList{Id: 2, Name: "B"}

	resp := &finpb.ListOperationsResponse{
		Operations: []*finpb.OperationInList{op1, op2},
	}

	res := MapListOperationsResponse(resp)
	require.Len(t, res, 2)
	require.Equal(t, "A", res[0].Name)
	require.Equal(t, "B", res[1].Name)
}

func TestMapListOperationsResponse_Nil(t *testing.T) {
	require.Nil(t, MapListOperationsResponse(nil))
}

func TestCreateOperationRequestToProto(t *testing.T) {
	date := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	catID := 5
	req := models.CreateOperationRequest{
		CategoryID:  &catID,
		Type:        "expense",
		Description: "Desc",
		Name:        "Op",
		Sum:         100.0,
		Date:        &date,
	}

	pb := CreateOperationRequestToProto(req, 10, 7)

	require.Equal(t, int32(10), pb.UserId)
	require.Equal(t, int32(7), pb.AccountId)
	require.Equal(t, int32(5), *pb.CategoryId)
	require.Equal(t, "expense", pb.Type)
	require.Equal(t, "Desc", pb.Description)
	require.Equal(t, "Op", pb.Name)
	require.Equal(t, 100.0, pb.Sum)
	require.Equal(t, date, pb.Date.AsTime())
}

func TestProtoOperationToResponse(t *testing.T) {
	created := timestamppb.New(time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC))
	date := timestamppb.New(time.Date(2024, 2, 2, 11, 0, 0, 0, time.UTC))

	op := &finpb.Operation{
		Id:           1,
		AccountId:    2,
		CategoryId:   3,
		CategoryName: "Food",
		Type:         "expense",
		Status:       "done",
		Description:  "Desc",
		ReceiptUrl:   "url",
		Name:         "Op",
		Sum:          123.4,
		CurrencyId:   840,
		CreatedAt:    created,
		Date:         date,
	}

	res := ProtoOperationToResponse(op)
	require.Equal(t, 1, res.ID)
	require.Equal(t, 2, res.AccountID)
	require.Equal(t, 3, res.CategoryID)
	require.Equal(t, "Food", res.CategoryName)
	require.Equal(t, "expense", res.Type)
	require.Equal(t, "done", res.Status)
	require.Equal(t, "Desc", res.Description)
	require.Equal(t, "url", res.ReceiptURL)
	require.Equal(t, "Op", res.Name)
	require.Equal(t, 123.4, res.Sum)
	require.Equal(t, 840, res.CurrencyID)
	require.Equal(t, created.AsTime(), res.CreatedAt)
	require.Equal(t, date.AsTime(), res.Date)
}

func TestOperationAndUserIDToProtoID(t *testing.T) {
	res := OperationAndUserIDToProtoID(3, 7, 10)
	require.Equal(t, int32(3), res.OperationId)
	require.Equal(t, int32(7), res.AccountId)
	require.Equal(t, int32(10), res.UserId)
}

func TestUpdateOperationRequestToProto(t *testing.T) {
	date := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
	descr := "descr"
	name := "name"
	sum := 200.0
	catID := 5
	req := models.UpdateOperationRequest{
		CategoryID:  &catID,
		Description: &descr,
		Name:        &name,
		Sum:         &sum,
		CreatedAt:   &date,
	}

	pb := UpdateOperationRequestToProto(req, 10, 7, 3)

	require.Equal(t, int32(10), pb.UserId)
	require.Equal(t, int32(7), pb.AccountId)
	require.Equal(t, int32(3), pb.OperationId)
	require.Equal(t, int32(5), *pb.CategoryId)
	require.Equal(t, "descr", *pb.Description)
	require.Equal(t, "name", *pb.Name)
	require.Equal(t, 200.0, *pb.Sum)
	require.Equal(t, date, pb.CreatedAt.AsTime())
}
