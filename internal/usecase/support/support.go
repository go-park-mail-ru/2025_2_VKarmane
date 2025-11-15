package support

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	supportsvc "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/support"
)

type UseCase struct {
	svc supportsvc.SupportService
}

func NewUseCase(svc supportsvc.SupportService) *UseCase {
	return &UseCase{
		svc: svc,
	}
}

func (uc *UseCase) CreateSupportRequest(ctx context.Context, userID int, category models.CategoryContacting, msg string) (models.Support, error) {
	log := logger.FromContext(ctx)

	supportReq, err := uc.svc.Create(ctx, userID, category, msg)
	if err != nil {
		log.Error("Failed to create support request",
			"error", err,
			"user_id", userID,
			"category", category,
		)
		return models.Support{}, pkgerrors.Wrap(err, "support.CreateSupportRequest")
	}

	return supportReq, nil
}

func (uc *UseCase) GetUserSupportRequests(ctx context.Context, userID int) ([]models.Support, error) {
	log := logger.FromContext(ctx)

	reqs, err := uc.svc.GetUserRequests(ctx, userID)
	if err != nil {
		log.Error("Failed to get support requests",
			"error", err,
			"user_id", userID,
		)
		return nil, pkgerrors.Wrap(err, "support.GetUserSupportRequests")
	}

	return reqs, nil
}

func (uc *UseCase) GetSupportStats(ctx context.Context) (map[models.StatusContacting]int, error) {
	log := logger.FromContext(ctx)

	stats, err := uc.svc.GetStats(ctx)
	if err != nil {
		log.Error("Failed to get support stats",
			"error", err,
		)
		return nil, pkgerrors.Wrap(err, "support.GetSupportStats")
	}

	return stats, nil
}

func (uc *UseCase) UpdateSupportStatus(ctx context.Context, id int, status models.StatusContacting) error {
	log := logger.FromContext(ctx)

	err := uc.svc.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error("Failed to update support status",
			"error", err,
			"support_id", id,
			"new_status", status,
		)
		return pkgerrors.Wrap(err, "support.UpdateSupportStatus")
	}

	return nil
}
