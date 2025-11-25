package image

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository"
)

type Service struct {
	storage image.ImageStorage
}

func NewService(storage image.ImageStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UploadImage(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	hash := sha256.Sum256(data)
	imageID := fmt.Sprintf("%x", hash)
	objectName := fmt.Sprintf("%s/%s", imageID[:2], imageID)

	_, err = s.storage.UploadImage(ctx, bytes.NewReader(data), objectName, int64(len(data)), contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to storage: %w", err)
	}

	return imageID, nil
}

func (s *Service) GetImageURL(ctx context.Context, imageID string) (string, error) {
	objectName := fmt.Sprintf("%s/%s", imageID[:2], imageID)
	url, err := s.storage.GetImageURL(ctx, objectName)
	if err != nil {
		return "", fmt.Errorf("failed to get image URL: %w", err)
	}

	return url, nil
}

func (s *Service) DeleteImage(ctx context.Context, imageID string) error {
	objectName := fmt.Sprintf("%s/%s", imageID[:2], imageID)
	err := s.storage.DeleteImage(ctx, objectName)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (s *Service) ImageExists(ctx context.Context, imageID string) (bool, error) {
	objectName := fmt.Sprintf("%s/%s", imageID[:2], imageID)
	exists, err := s.storage.ImageExists(ctx, objectName)
	if err != nil {
		return false, fmt.Errorf("failed to check image existence: %w", err)
	}

	return exists, nil
}
