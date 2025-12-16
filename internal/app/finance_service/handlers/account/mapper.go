package account

import (
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type AccountAPI struct {
	ID          int     `json:"id"`
	Balance     float64 `json:"balance"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CurrencyID  int     `json:"currency_id"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type AccountsAPI struct {
	UserID   int          `json:"user_id"`
	Accounts []AccountAPI `json:"accounts"`
	TotalSum float64      `json:"total_sum"`
}

type SharingApi struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	AccountID int    `json:"account_id"`
	CreatedAt string `json:"created_at"`
}

func UserIDToProtoID(userID int) *finpb.UserID {
	return &finpb.UserID{
		UserId: int32(userID),
	}
}

func UserIDAndAccountIDToProtoID(userID, accID int) *finpb.AccountRequest {
	return &finpb.AccountRequest{
		UserId:    int32(userID),
		AccountId: int32(accID),
	}
}

func UserLoginIDtoProtoID(login string, accID int) *finpb.AddToAccountReqeust {
	return &finpb.AddToAccountReqeust{
		AccountId: int32(accID),
		UserLogin: login,
	}
}

func AccountResponseListProtoToApi(resp *finpb.ListAccountsResponse, userID int) AccountsAPI {
	if resp == nil {
		return AccountsAPI{
			UserID:   userID,
			Accounts: make([]AccountAPI, 0),
		}
	}

	var sum float64
	accounts := make([]AccountAPI, 0, len(resp.Accounts))

	for _, acc := range resp.Accounts {
		if acc == nil {
			continue
		}

		var createdAt string
		if acc.CreatedAt != nil {
			createdAt = acc.CreatedAt.AsTime().Format(time.RFC3339)
		}

		var updatedAt string
		if acc.UpdatedAt != nil {
			updatedAt = acc.UpdatedAt.AsTime().Format(time.RFC3339)
		}

		accounts = append(accounts, AccountAPI{
			ID:          int(acc.Id),
			Balance:     acc.Balance,
			Name:        acc.Name,
			Description: acc.Description,
			Type:        acc.Type,
			CurrencyID:  int(acc.CurrencyId),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
		sum += acc.Balance
	}

	return AccountsAPI{
		UserID:   userID,
		Accounts: accounts,
		TotalSum: sum,
	}
}

func ProtoAccountToApi(acc *finpb.Account) AccountAPI {
	return AccountAPI{
		ID:          int(acc.Id),
		Balance:     acc.Balance,
		Name:        acc.Name,
		Description: acc.Description,
		Type:        acc.Type,
		CurrencyID:  int(acc.CurrencyId),
		CreatedAt:   acc.CreatedAt.AsTime().Format(time.RFC3339),
		UpdatedAt:   acc.UpdatedAt.AsTime().Format(time.RFC3339),
	}
}

func AccountCreateRequestToProto(userID int, req models.CreateAccountRequest) *finpb.CreateAccountRequest {
	return &finpb.CreateAccountRequest{
		UserId:      int32(userID),
		Balance:     req.Balance,
		CurrencyId:  int32(req.CurrencyID),
		Type:        string(req.Type),
		Name:        req.Name,
		Description: req.Description,
	}
}

func AccountUpdateRequestToProto(userID, accID int, req models.UpdateAccountRequest) *finpb.UpdateAccountRequest {
	return &finpb.UpdateAccountRequest{
		UserId:      int32(userID),
		AccountId:   int32(accID),
		Balance:     req.Balance,
		Name:        req.Name,
		Description: req.Description,
	}
}

func SharingProtoToApi(resp *finpb.SharingsResponse) SharingApi {
	return SharingApi{
		ID:        int(resp.SharingId),
		AccountID: int(resp.AccountId),
		UserID:    int(resp.UserId),
		CreatedAt: resp.CreatedAt.AsTime().Format(time.RFC3339),
	}
}
