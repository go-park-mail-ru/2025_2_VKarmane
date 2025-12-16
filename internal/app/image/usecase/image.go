package image

import (
	"context"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

type UseCase struct {
	imageSvc image.ImageService
}

func NewUseCase(imageSvc image.ImageService) *UseCase {
	return &UseCase{
		imageSvc: imageSvc,
	}
}

func (uc *UseCase) UploadImage(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	log := logger.FromContext(ctx)
	imageID, err := uc.imageSvc.UploadImage(ctx, reader, filename, size, contentType)
	if err != nil {
		log.Error("Failed to upload image", "error", err, "filename", filename)
		return "", fmt.Errorf("image.UploadImage: %w", err)
	}

	return imageID, nil
}

func (uc *UseCase) GetImageURL(ctx context.Context, imageID string) (string, error) {
	url, err := uc.imageSvc.GetImageURL(ctx, imageID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get image URL", "error", err, "image_id", imageID)
		}
		return "", fmt.Errorf("image.GetImageURL: %w", err)
	}

	return url, nil
}

func (uc *UseCase) DeleteImage(ctx context.Context, imageID string) error {
	log := logger.FromContext(ctx)
	err := uc.imageSvc.DeleteImage(ctx, imageID)
	if err != nil {
		log.Error("Failed to delete image", "error", err, "image_id", imageID)
		return fmt.Errorf("image.DeleteImage: %w", err)
	}

	return nil
}

func (uc *UseCase) ImageExists(ctx context.Context, imageID string) (bool, error) {
	log := logger.FromContext(ctx)
	exists, err := uc.imageSvc.ImageExists(ctx, imageID)
	if err != nil {
		log.Error("Failed to check image existence", "error", err, "image_id", imageID)
		return false, fmt.Errorf("image.ImageExists: %w", err)
	}

	return exists, nil
}
