package image

import (
	"context"
	"io"
)

type ImageService interface {
	UploadImage(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error)
	GetImageURL(ctx context.Context, imageID string) (string, error)
	DeleteImage(ctx context.Context, imageID string) error
	ImageExists(ctx context.Context, imageID string) (bool, error)
}
