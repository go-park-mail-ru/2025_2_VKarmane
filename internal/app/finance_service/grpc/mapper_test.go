package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
)

func TestProtoToCreateAccountRequest(t *testing.T) {
	req := &finpb.CreateAccountRequest{UserId: 3, Balance: 100, Type: "cash", CurrencyId: 1}
	result := protoToCreateAccountRequest(req)

	require.Equal(t, finmodels.CreateAccountRequest{
		UserID:     3,
		Balance:    100,
		Type:       "cash",
		CurrencyID: 1,
	}, result)
}

func TestProtoToCreateOperationRequest(t *testing.T) {
	now := time.Now()
	categoryID := int32(5)
	req := &finpb.CreateOperationRequest{
		UserId:      1,
		AccountId:   2,
		CategoryId:  &categoryID,
		Type:        "expense",
		Name:        "groceries",
		Description: "food",
		Sum:         50,
		Date:        timestamppb.New(now),
	}

	result := protoToCreateOperationRequest(req)
	require.Equal(t, 1, result.UserID)
	require.Equal(t, 2, result.AccountID)
	require.NotNil(t, result.CategoryID)
	require.Equal(t, 5, *result.CategoryID)
	require.NotNil(t, result.Date)
	require.WithinDuration(t, now, result.Date.UTC(), time.Millisecond)
}

func TestProtoToUpdateOperationRequestOptional(t *testing.T) {
	categoryID := int32(7)
	name := "rent"
	description := "march"
	sum := 400.0
	now := time.Now()
	req := &finpb.UpdateOperationRequest{
		UserId:      1,
		AccountId:   2,
		OperationId: 3,
		CategoryId:  &categoryID,
		Name:        &name,
		Description: &description,
		Sum:         &sum,
		CreatedAt:   timestamppb.New(now),
	}

	result := protoToUpdateOperationRequest(req)
	require.NotNil(t, result.CategoryID)
	require.Equal(t, 7, *result.CategoryID)
	require.NotNil(t, result.Name)
	require.Equal(t, name, *result.Name)
	require.NotNil(t, result.CreatedAt)
	require.WithinDuration(t, now, result.CreatedAt.UTC(), time.Millisecond)
}

func TestProtoToUpdateCategoryRequest(t *testing.T) {
	name := "updated"
	description := "desc"
	logo := "logo"
	req := &finpb.UpdateCategoryRequest{
		UserId:       1,
		CategoryId:   2,
		Name:         &name,
		Description:  &description,
		LogoHashedId: &logo,
	}

	result := protoToUpdateCategoryRequest(req)
	require.Equal(t, 1, result.UserID)
	require.Equal(t, 2, result.CategoryID)
	require.Equal(t, &name, result.Name)
	require.Equal(t, &logo, result.LogoHashedID)
}