package image

import (
	"context"
	"io"
)

type ImageStorage interface {
	UploadImage(ctx context.Context, reader io.Reader, objectName string, size int64, contentType string) (string, error)
	GetImageURL(ctx context.Context, objectName string) (string, error)
	DeleteImage(ctx context.Context, objectName string) error
	ImageExists(ctx context.Context, objectName string) (bool, error)
}
