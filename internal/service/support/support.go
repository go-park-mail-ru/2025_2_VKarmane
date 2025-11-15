package support

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"
)

type repo interface {
	Create(ctx context.Context, s models.Support) (models.Support, error)
	GetByUser(ctx context.Context, userID int) ([]models.Support, error)
	GetStats(ctx context.Context) (map[models.StatusContacting]int, error)
	UpdateStatus(ctx context.Context, id int, status models.StatusContacting) error
}

type Service struct {
	repo repo
}

func NewService(r repo) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CheckOwnership(ctx context.Context, userID int) bool {
	uid, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || uid == 0 {
		return false
	}
	return uid == userID
}

func (s *Service) Create(ctx context.Context, userID int, category models.CategoryContacting, msg string) (models.Support, error) {
	if !s.CheckOwnership(ctx, userID) {
		return models.Support{}, serviceerrors.ErrForbidden
	}

	if msg == "" {
		return models.Support{}, pkgerrors.Wrap(serviceerrors.ErrInvalidInput, "empty message")
	}

	model := models.Support{
		UserID:          userID,
		CategoryRequest: category,
		StatusRequest:   models.RequestOpened,
		Message:         msg,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return models.Support{}, pkgerrors.Wrap(err, "support.Create")
	}

	return created, nil
}

func (s *Service) GetUserRequests(ctx context.Context, userID int) ([]models.Support, error) {
	if !s.CheckOwnership(ctx, userID) {
		return nil, serviceerrors.ErrForbidden
	}

	list, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return []models.Support{}, pkgerrors.Wrap(err, "support.GetUserRequests")
	}

	return list, nil
}

func (s *Service) GetStats(ctx context.Context) (map[models.StatusContacting]int, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "support.GetStats")
	}

	return stats, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status models.StatusContacting) error {
	// Проверка прав (например, только админ)
	//isAdmin, _ := middleware.GetIsAdminFromContext(ctx)
	//if !isAdmin {
	//	return serviceerrors.ErrForbidden
	//}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return pkgerrors.Wrap(err, "support.UpdateStatus")
	}

	return nil
}
