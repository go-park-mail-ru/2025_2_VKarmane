package category

import (
	"net/url"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func CreateCategoryReportRequestToProto(userID int, q url.Values) *finpb.CategoryReportRequest {
	startStr := q.Get("start")
	endStr := q.Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			start = time.Time{}
		}
	}

	if endStr != "" {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			end = time.Time{}
		}
	}

	return &finpb.CategoryReportRequest{
		UserId: int32(userID),
		Start:  timestamppb.New(start),
		End:    timestamppb.New(end),
	}
}

func ProtoToCategoryReport(resp *finpb.CategoryReportResponse) models.CategoryReport {
	result := models.CategoryReport{
		Categories: make([]models.CategoryInReport, 0, len(resp.GetCategories())),
		Start:      resp.GetStart().AsTime(),
		End:        resp.GetEnd().AsTime(),
	}

	for _, c := range resp.GetCategories() {
		result.Categories = append(result.Categories, models.CategoryInReport{
			CategoryID:     int(c.CategoryId),
			CategoryName:   c.CategoryName,
			OperationCount: int(c.OperationsCount),
			TotalSum:       c.TotalSum,
		})
	}

	return result
}
