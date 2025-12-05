package category

import (
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func UserIDToProtoID(userID int) *finpb.UserID {
	return &finpb.UserID{
		UserId: int32(userID),
	}
}

func CategoryWithStatsToAPI(ctg *finpb.CategoryWithStats) models.CategoryWithStats {
	return models.CategoryWithStats{
		Category: models.Category{
			ID:           int(ctg.Category.Id),
			UserID:       int(ctg.Category.UserId),
			Name:         ctg.Category.Name,
			Description:  ctg.Category.Description,
			LogoHashedID: ctg.Category.LogoHashedId,
			LogoURL:      ctg.Category.LogoUrl,
			CreatedAt:    ctg.Category.CreatedAt.AsTime(),
			UpdatedAt:    ctg.Category.UpdatedAt.AsTime(),
		},
		OperationsCount: int(ctg.OperationsCount),
	}
}
func CategoriesWithStatsToAPI(userID int, ctgs *finpb.ListCategoriesWithStatsResponse) []models.CategoryWithStats {
	res := make([]models.CategoryWithStats, 0, len(ctgs.Categories))
	for _, b := range ctgs.Categories {
		res = append(res, CategoryWithStatsToAPI(b))
	}
	return res
}

func CategoryCreateRequestToProto(userID int, req models.CreateCategoryRequest) *finpb.CreateCategoryRequest {
	return &finpb.CreateCategoryRequest{
		UserId:       int32(userID),
		Name:         req.Name,
		Description:  req.Description,
		LogoHashedId: req.LogoHashedID,
	}
}

func ProtoCategoryToApi(ctg *finpb.Category) *models.Category {
	return &models.Category{
		ID:           int(ctg.Id),
		UserID:       int(ctg.UserId),
		Name:         ctg.Name,
		Description:  ctg.Description,
		LogoHashedID: ctg.LogoHashedId,
		LogoURL:      ctg.LogoUrl,
		CreatedAt:    ctg.CreatedAt.AsTime(),
		UpdatedAt:    ctg.UpdatedAt.AsTime(),
	}
}

func UserAndCtegoryIDToProto(userID, ctgID int) *finpb.CategoryRequest {
	return &finpb.CategoryRequest{
		UserId:     int32(userID),
		CategoryId: int32(ctgID),
	}
}

func UserIDCategoryNameToProto(userID int, categoryName string) *finpb.CategoryByNameRequest {
	return &finpb.CategoryByNameRequest{
		UserId:       int32(userID),
		CategoryName: categoryName,
	}
}

func CategoryUpdateRequestToProto(userID, ctgID int, req models.UpdateCategoryRequest) *finpb.UpdateCategoryRequest {
	return &finpb.UpdateCategoryRequest{
		UserId:       int32(userID),
		CategoryId:   int32(ctgID),
		Name:         req.Name,
		Description:  req.Description,
		LogoHashedId: req.LogoHashedID,
	}
}

func CategoryToUpdateSearch(ctg *models.Category) models.UpdateCategoryInOperationSearch {
	return models.UpdateCategoryInOperationSearch{
		CategoryID:           ctg.ID,
		CategoryName:         ctg.Name,
		CategoryLogo:         ctg.LogoURL,
		CategoryLogoHashedID: ctg.LogoHashedID,
	}
}
