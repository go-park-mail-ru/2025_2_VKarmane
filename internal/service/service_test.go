package service

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	imagerepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository"
)

type stubStorage struct{}

func (stubStorage) UploadImage(ctx context.Context, reader io.Reader, objectName string, size int64, contentType string) (string, error) {
	return objectName, nil
}

func (stubStorage) GetImageURL(ctx context.Context, objectName string) (string, error) {
	return "http://cdn/" + objectName, nil
}

func (stubStorage) DeleteImage(ctx context.Context, objectName string) error {
	return nil
}

func (stubStorage) ImageExists(ctx context.Context, objectName string) (bool, error) {
	return true, nil
}

func TestNewService(t *testing.T) {
	service := NewService("secret", stubStorage{})
	require.NotNil(t, service)
	require.NotNil(t, service.ImageUC)
}

var _ imagerepo.ImageStorage = (*stubStorage)(nil)