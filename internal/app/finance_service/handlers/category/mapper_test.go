package category

import (
	"testing"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUserIDToProtoID(t *testing.T) {
	res := UserIDToProtoID(15)
	require.Equal(t, int32(15), res.UserId)
}

func TestCategoryWithStatsToAPI(t *testing.T) {
	created := timestamppb.New(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
	updated := timestamppb.New(time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC))

	ctg := &finpb.CategoryWithStats{
		Category: &finpb.Category{
			Id:           1,
			UserId:       42,
			Name:         "Food",
			Description:  "Food expenses",
			LogoHashedId: "hash123",
			LogoUrl:      "http://logo.url",
			CreatedAt:    created,
			UpdatedAt:    updated,
		},
		OperationsCount: 5,
	}

	api := CategoryWithStatsToAPI(ctg)

	require.Equal(t, 1, api.Category.ID)
	require.Equal(t, 42, api.Category.UserID)
	require.Equal(t, "Food", api.Category.Name)
	require.Equal(t, "Food expenses", api.Category.Description)
	require.Equal(t, "hash123", api.Category.LogoHashedID)
	require.Equal(t, "http://logo.url", api.Category.LogoURL)
	require.Equal(t, created.AsTime(), api.Category.CreatedAt)
	require.Equal(t, updated.AsTime(), api.Category.UpdatedAt)
	require.Equal(t, 5, api.OperationsCount)
}

func TestCategoriesWithStatsToAPI(t *testing.T) {
	created := timestamppb.Now()
	updated := timestamppb.Now()

	resp := &finpb.ListCategoriesWithStatsResponse{
		Categories: []*finpb.CategoryWithStats{
			{
				Category: &finpb.Category{
					Id:           1,
					UserId:       1,
					Name:         "A",
					Description:  "Desc",
					LogoHashedId: "h1",
					LogoUrl:      "url1",
					CreatedAt:    created,
					UpdatedAt:    updated,
				},
				OperationsCount: 10,
			},
			{
				Category: &finpb.Category{
					Id:           2,
					UserId:       1,
					Name:         "B",
					Description:  "Desc2",
					LogoHashedId: "h2",
					LogoUrl:      "url2",
					CreatedAt:    created,
					UpdatedAt:    updated,
				},
				OperationsCount: 20,
			},
		},
	}

	api := CategoriesWithStatsToAPI(1, resp)
	require.Len(t, api, 2)
	require.Equal(t, "A", api[0].Category.Name)
	require.Equal(t, "B", api[1].Category.Name)
	require.Equal(t, 10, api[0].OperationsCount)
	require.Equal(t, 20, api[1].OperationsCount)
}

func TestCategoryCreateRequestToProto(t *testing.T) {
	req := models.CreateCategoryRequest{
		Name:         "Test",
		Description:  "Desc",
		LogoHashedID: "hashX",
	}

	pb := CategoryCreateRequestToProto(5, req)

	require.Equal(t, int32(5), pb.UserId)
	require.Equal(t, "Test", pb.Name)
	require.Equal(t, "Desc", pb.Description)
	require.Equal(t, "hashX", pb.LogoHashedId)
}

func TestProtoCategoryToApi(t *testing.T) {
	created := timestamppb.New(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	updated := timestamppb.New(time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC))

	ctg := &finpb.Category{
		Id:           3,
		UserId:       10,
		Name:         "TestCtg",
		Description:  "DescCtg",
		LogoHashedId: "hashC",
		LogoUrl:      "urlC",
		CreatedAt:    created,
		UpdatedAt:    updated,
	}

	api := ProtoCategoryToApi(ctg)
	require.Equal(t, 3, api.ID)
	require.Equal(t, 10, api.UserID)
	require.Equal(t, "TestCtg", api.Name)
	require.Equal(t, "DescCtg", api.Description)
	require.Equal(t, "hashC", api.LogoHashedID)
	require.Equal(t, "urlC", api.LogoURL)
	require.Equal(t, created.AsTime(), api.CreatedAt)
	require.Equal(t, updated.AsTime(), api.UpdatedAt)
}

func TestUserAndCtegoryIDToProto(t *testing.T) {
	pb := UserAndCtegoryIDToProto(5, 7)
	require.Equal(t, int32(5), pb.UserId)
	require.Equal(t, int32(7), pb.CategoryId)
}

func TestCategoryUpdateRequestToProto(t *testing.T) {
	name := "updated"
	descr := "descr"
	logo := "hash"
	req := models.UpdateCategoryRequest{
		Name:         &name,
		Description:  &descr,
		LogoHashedID: &logo,
	}

	pb := CategoryUpdateRequestToProto(1, 2, req)

	require.Equal(t, int32(1), pb.UserId)
	require.Equal(t, int32(2), pb.CategoryId)
	require.Equal(t, "updated", *pb.Name)
	require.Equal(t, "descr", *pb.Description)
	require.Equal(t, "hash", *pb.LogoHashedId)
}
