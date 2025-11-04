package image

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_UploadImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockImageService(ctrl)
	uc := NewUseCase(mockSvc)

	file := bytes.NewReader([]byte("test"))

	mockSvc.EXPECT().UploadImage(gomock.Any(), file, "test.jpg", int64(4), "image/jpeg").Return("img-123", nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.UploadImage(ctx, file, "test.jpg", 4, "image/jpeg")
	assert.NoError(t, err)
	assert.Equal(t, "img-123", result)
}

func TestUseCase_UploadImage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockImageService(ctrl)
	uc := NewUseCase(mockSvc)

	file := bytes.NewReader([]byte("test"))

	mockSvc.EXPECT().UploadImage(gomock.Any(), file, "test.jpg", int64(4), "image/jpeg").Return("", errors.New("upload failed"))

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.UploadImage(ctx, file, "test.jpg", 4, "image/jpeg")
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestUseCase_GetImageURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockImageService(ctrl)
	uc := NewUseCase(mockSvc)

	mockSvc.EXPECT().GetImageURL(gomock.Any(), "img-123").Return("http://example.com/img.jpg", nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetImageURL(ctx, "img-123")
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/img.jpg", result)
}

func TestUseCase_DeleteImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockImageService(ctrl)
	uc := NewUseCase(mockSvc)

	mockSvc.EXPECT().DeleteImage(gomock.Any(), "img-123").Return(nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	err := uc.DeleteImage(ctx, "img-123")
	assert.NoError(t, err)
}
