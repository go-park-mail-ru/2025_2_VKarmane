package image

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_UploadImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockImageStorage(ctrl)
	svc := NewService(mockStorage)

	file := bytes.NewReader([]byte("test image data"))
	
	mockStorage.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("stored-hash", nil)

	result, err := svc.UploadImage(context.Background(), file, "test.jpg", 15, "image/jpeg")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestService_GetImageURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockImageStorage(ctrl)
	svc := NewService(mockStorage)

	mockStorage.EXPECT().GetImageURL(gomock.Any(), gomock.Any()).Return("http://example.com/image.jpg", nil)

	result, err := svc.GetImageURL(context.Background(), "image-id-123")
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/image.jpg", result)
}

func TestService_DeleteImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockImageStorage(ctrl)
	svc := NewService(mockStorage)

	mockStorage.EXPECT().DeleteImage(gomock.Any(), gomock.Any()).Return(nil)

	err := svc.DeleteImage(context.Background(), "image-id-123")
	assert.NoError(t, err)
}

