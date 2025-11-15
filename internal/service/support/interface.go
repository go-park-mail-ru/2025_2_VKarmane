package support

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type SupportService interface {
	Create(ctx context.Context, userID int, category models.CategoryContacting, msg string) (models.Support, error)
	GetUserRequests(ctx context.Context, userID int) ([]models.Support, error)
	GetStats(ctx context.Context) (map[models.StatusContacting]int, error)
	UpdateStatus(ctx context.Context, id int, status models.StatusContacting) error
}
