package account

import (
	"testing"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUserIDToProtoID(t *testing.T) {
	res := UserIDToProtoID(10)
	require.Equal(t, int32(10), res.UserId)
}

func TestUserIDAndAccountIDToProtoID(t *testing.T) {
	res := UserIDAndAccountIDToProtoID(10, 7)
	require.Equal(t, int32(10), res.UserId)
	require.Equal(t, int32(7), res.AccountId)
}

func TestAccountResponseListProtoToApit_NilResponse(t *testing.T) {
	res := AccountResponseListProtoToApit(nil, 42)

	require.Equal(t, 42, res.UserID)
	require.Len(t, res.Accounts, 0)
}

func TestAccountResponseListProtoToApit_Full(t *testing.T) {
	t1 := timestamppb.New(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	t2 := timestamppb.New(time.Date(2024, 1, 2, 13, 0, 0, 0, time.UTC))

	resp := &finpb.ListAccountsResponse{
		Accounts: []*finpb.Account{
			{
				Id:         1,
				Balance:    100.5,
				Type:       "cash",
				CurrencyId: 643,
				CreatedAt:  t1,
				UpdatedAt:  t2,
			},
		},
	}

	api := AccountResponseListProtoToApit(resp, 99)

	require.Equal(t, 99, api.UserID)
	require.Len(t, api.Accounts, 1)

	acc := api.Accounts[0]
	require.Equal(t, 1, acc.ID)
	require.Equal(t, 100.5, acc.Balance)
	require.Equal(t, "cash", acc.Type)
	require.Equal(t, 643, acc.CurrencyID)
	require.Equal(t, t1.AsTime().Format(time.RFC3339), acc.CreatedAt)
	require.Equal(t, t2.AsTime().Format(time.RFC3339), acc.UpdatedAt)
}

func TestProtoAccountToApi(t *testing.T) {
	t1 := timestamppb.New(time.Date(2024, 5, 1, 10, 20, 30, 0, time.UTC))
	t2 := timestamppb.New(time.Date(2024, 5, 2, 11, 22, 33, 0, time.UTC))

	pb := &finpb.Account{
		Id:         3,
		Balance:    200.75,
		Type:       "card",
		CurrencyId: 840,
		CreatedAt:  t1,
		UpdatedAt:  t2,
	}

	api := ProtoAccountToApi(pb)

	require.Equal(t, 3, api.ID)
	require.Equal(t, 200.75, api.Balance)
	require.Equal(t, "card", api.Type)
	require.Equal(t, 840, api.CurrencyID)
	require.Equal(t, t1.AsTime().Format(time.RFC3339), api.CreatedAt)
	require.Equal(t, t2.AsTime().Format(time.RFC3339), api.UpdatedAt)
}

func TestAccountCreateRequestToProto(t *testing.T) {
	req := models.CreateAccountRequest{
		Balance:    150.5,
		CurrencyID: 978,
		Type:       "deposit",
	}

	pb := AccountCreateRequestToProto(12, req)

	require.Equal(t, int32(12), pb.UserId)
	require.Equal(t, 150.5, pb.Balance)
	require.Equal(t, int32(978), pb.CurrencyId)
	require.Equal(t, "deposit", pb.Type)
}

func TestAccountUpdateRequestToProto(t *testing.T) {
	req := models.UpdateAccountRequest{
		Balance: 300.99,
	}

	pb := AccountUpdateRequestToProto(5, 11, req)

	require.Equal(t, int32(5), pb.UserId)
	require.Equal(t, int32(11), pb.AccountId)
	require.Equal(t, 300.99, pb.Balance)
}
