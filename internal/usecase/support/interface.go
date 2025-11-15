package support

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type SupportService interface {
	CreateSupportRequest(ctx context.Context, userID int, category models.CategoryContacting, msg string) (models.Support, error)
	GetUserSupportRequests(ctx context.Context, userID int) ([]models.Support, error)
	GetSupportStats(ctx context.Context) (map[models.StatusContacting]int, error)
	UpdateSupportStatus(ctx context.Context, id int, status models.StatusContacting) error
}
