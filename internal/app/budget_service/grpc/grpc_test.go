package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/errors"
	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
)

func TestBudgetServiceServer_CreateBudgetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockBudgetUseCase(ctrl)
	server := NewBudgetServer(uc)
	req := &bdgpb.CreateBudgetRequest{UserID: 1, CategoryId: 2, Sum: 10}
	expected := &bdgpb.Budget{Id: 1, UserId: 1}

	uc.EXPECT().
		CreateBudget(gomock.Any(), gomock.AssignableToTypeOf(bdgmodels.CreateBudgetRequest{}), 1).
		Return(expected, nil)

	resp, err := server.CreateBudget(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestBudgetServiceServer_CreateBudgetKnownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockBudgetUseCase(ctrl)
	server := NewBudgetServer(uc)
	req := &bdgpb.CreateBudgetRequest{UserID: 1}

	uc.EXPECT().
		CreateBudget(gomock.Any(), gomock.AssignableToTypeOf(bdgmodels.CreateBudgetRequest{}), 1).
		Return(nil, bdgerrors.ErrBudgetExists)

	resp, err := server.CreateBudget(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.AlreadyExists, st.Code())
}

func TestBudgetServiceServer_GetBudgetInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockBudgetUseCase(ctrl)
	server := NewBudgetServer(uc)
	req := &bdgpb.BudgetRequest{BudgetID: 1, UserID: 1}

	uc.EXPECT().GetBudget(gomock.Any(), 1, 1).Return(nil, errors.New("boom"))

	resp, err := server.GetBudgetByID(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Internal, st.Code())
}

func TestBudgetServiceServer_UpdateBudgetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockBudgetUseCase(ctrl)
	server := NewBudgetServer(uc)
	sum := 50.0
	req := &bdgpb.UpdateBudgetRequest{BudgetID: 1, UserID: 1, Sum: &sum}
	expected := &bdgpb.Budget{Id: 1, UserId: 1, Sum: 50}

	uc.EXPECT().
		UpdateBudget(gomock.Any(), gomock.AssignableToTypeOf(bdgmodels.UpdatedBudgetRequest{})).
		Return(expected, nil)

	resp, err := server.UpdateBudget(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestBudgetServiceServer_DeleteBudgetKnownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockBudgetUseCase(ctrl)
	server := NewBudgetServer(uc)
	req := &bdgpb.BudgetRequest{BudgetID: 1, UserID: 1}

	uc.EXPECT().DeleteBudget(gomock.Any(), 1, 1).Return(nil, bdgerrors.ErrBudgetNotFound)

	resp, err := server.DeleteBudget(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, st.Code())
}

func TestProtoMappers(t *testing.T) {
	now := time.Now()
	sum := 50.0
	req := &bdgpb.CreateBudgetRequest{
		UserID:      1,
		CategoryId:  2,
		Sum:         100,
		Description: "test",
		CreatedAt:   timestamppb.New(now),
		PeriodStart: timestamppb.New(now),
		PeriodEnd:   timestamppb.New(now.Add(time.Hour)),
	}
	model, userID := ProtoCreateRequestToModel(req)
	require.Equal(t, 1, userID)
	require.Equal(t, 2, model.CategoryID)

	description := "upd"
	updateReq := &bdgpb.UpdateBudgetRequest{
		UserID:      1,
		BudgetID:    2,
		Sum:         &sum,
		Description: &description,
		PeriodStart: timestamppb.New(now),
	}
	updateModel := ProtoUpdateRequestToModel(updateReq)
	require.Equal(t, 1, updateModel.UserID)
	require.Equal(t, 2, updateModel.BudgetID)
	require.NotNil(t, updateModel.PeriodStart)

	reqID := &bdgpb.UserID{UserID: 7}
	require.Equal(t, 7, ProtoIDToInt(reqID))

	breq := &bdgpb.BudgetRequest{UserID: 3, BudgetID: 4}
	user, budget := ProtoBudgetReqToInts(breq)
	require.Equal(t, 3, user)
	require.Equal(t, 4, budget)
}
