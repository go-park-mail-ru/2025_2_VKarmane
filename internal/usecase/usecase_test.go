package usecase

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	imagerepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository"
	imageservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
)

type fakeStorage struct{}

func (fakeStorage) UploadImage(ctx context.Context, reader io.Reader, objectName string, size int64, contentType string) (string, error) {
	return objectName, nil
}

func (fakeStorage) GetImageURL(ctx context.Context, objectName string) (string, error) {
	return "http://cdn/" + objectName, nil
}

func (fakeStorage) DeleteImage(ctx context.Context, objectName string) error {
	return nil
}

func (fakeStorage) ImageExists(ctx context.Context, objectName string) (bool, error) {
	return true, nil
}

func TestNewUseCase(t *testing.T) {
	svc := &service.Service{
		ImageUC: imageservice.NewService(fakeStorage{}),
	}

	uc := NewUseCase(svc, "secret")

	require.Equal(t, svc, uc.service)
	require.NotNil(t, uc.ImageUC)
}

var _ imagerepo.ImageStorage = (*fakeStorage)(nil)

