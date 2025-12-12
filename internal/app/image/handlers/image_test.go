package image

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
)

type fakeImageUseCase struct {
	uploadFunc func(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error)
	getURLFunc func(ctx context.Context, imageID string) (string, error)
	deleteFunc func(ctx context.Context, imageID string) error
	existsFunc func(ctx context.Context, imageID string) (bool, error)
}

func (f *fakeImageUseCase) UploadImage(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	if f.uploadFunc != nil {
		return f.uploadFunc(ctx, reader, filename, size, contentType)
	}
	return "image-id", nil
}

func (f *fakeImageUseCase) GetImageURL(ctx context.Context, imageID string) (string, error) {
	if f.getURLFunc != nil {
		return f.getURLFunc(ctx, imageID)
	}
	return "http://example.com/" + imageID, nil
}

func (f *fakeImageUseCase) DeleteImage(ctx context.Context, imageID string) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, imageID)
	}
	return nil
}

func (f *fakeImageUseCase) ImageExists(ctx context.Context, imageID string) (bool, error) {
	if f.existsFunc != nil {
		return f.existsFunc(ctx, imageID)
	}
	return true, nil
}

func withUser(req *http.Request, userID int) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func newMultipartRequest(t *testing.T, filename string, content []byte, includeContentType bool) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var part io.Writer
	var err error

	if includeContentType {
		part, err = writer.CreateFormFile("image", filename)
	} else {
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", filename))
		part, err = writer.CreatePart(header)
	}
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/images/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func decodeErrorResponse(t *testing.T, body *bytes.Buffer) models.ErrorResponse {
	t.Helper()
	var resp models.ErrorResponse
	require.NoError(t, json.Unmarshal(body.Bytes(), &resp))
	return resp
}

func TestUploadImage_Unauthorized(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := httptest.NewRequest(http.MethodPost, "/images/upload", nil)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, models.ErrCodeUnauthorized, resp.Code)
}

func TestUploadImage_InvalidMultipart(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := httptest.NewRequest(http.MethodPost, "/images/upload", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=test")
	req = withUser(req, 1)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "Failed to parse multipart form")
}

func TestUploadImage_InvalidExtension(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := newMultipartRequest(t, "note.txt", []byte("content"), true)
	req = withUser(req, 42)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "Invalid image format")
}

func TestUploadImage_MissingContentTypeDefaults(t *testing.T) {
	var capturedContentType string
	uc := &fakeImageUseCase{
		uploadFunc: func(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
			capturedContentType = contentType
			return "hash", nil
		},
	}
	handler := NewHandler(uc)
	req := newMultipartRequest(t, "photo.PNG", []byte("binary-data"), false)
	req = withUser(req, 7)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, "image/png", capturedContentType)
}

func TestUploadImage_UploadError(t *testing.T) {
	uc := &fakeImageUseCase{
		uploadFunc: func(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
			return "", errors.New("boom")
		},
	}
	handler := NewHandler(uc)
	req := newMultipartRequest(t, "photo.png", []byte("data"), true)
	req = withUser(req, 13)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, "Failed to upload image", resp.Error)
}

func TestUploadImage_GetURLError(t *testing.T) {
	uc := &fakeImageUseCase{
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			return "", errors.New("oops")
		},
	}
	handler := NewHandler(uc)
	req := newMultipartRequest(t, "photo.png", []byte("data"), true)
	req = withUser(req, 21)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, "Failed to get image URL", resp.Error)
}

func TestUploadImage_Success(t *testing.T) {
	uc := &fakeImageUseCase{
		uploadFunc: func(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
			require.Equal(t, "photo.png", filename)
			return "abc123", nil
		},
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			require.Equal(t, "abc123", imageID)
			return "http://cdn/images/abc123", nil
		},
	}
	handler := NewHandler(uc)
	req := newMultipartRequest(t, "photo.png", []byte("data"), true)
	req = withUser(req, 99)
	rr := httptest.NewRecorder()

	handler.UploadImage(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	var resp UploadImageResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, "abc123", resp.ImageID)
	require.Equal(t, "http://cdn/images/abc123", resp.URL)
}

func TestGetImageURL_Unauthorized(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := httptest.NewRequest(http.MethodGet, "/images/url", nil)
	rr := httptest.NewRecorder()

	handler.GetImageURL(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetImageURL_MissingID(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := withUser(httptest.NewRequest(http.MethodGet, "/images/url", nil), 1)
	rr := httptest.NewRecorder()

	handler.GetImageURL(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "image_id is required")
}

func TestGetImageURL_Error(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			return "", errors.New("boom")
		},
	})
	req := withUser(httptest.NewRequest(http.MethodGet, "/images/url?image_id=img", nil), 42)
	rr := httptest.NewRecorder()

	handler.GetImageURL(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, "Failed to get image URL", resp.Error)
}

func TestGetImageURL_Success(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			return "http://cdn/" + imageID, nil
		},
	})
	req := withUser(httptest.NewRequest(http.MethodGet, "/images/url?image_id=image123", nil), 5)
	rr := httptest.NewRecorder()

	handler.GetImageURL(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp ImageURLResponse
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, "http://cdn/image123", resp.URL)
}

func TestGetImage_MissingID(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := httptest.NewRequest(http.MethodGet, "/images", nil)
	rr := httptest.NewRecorder()

	handler.GetImage(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "image_id is required")
}

func TestGetImage_Error(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			return "", errors.New("boom")
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/images?image_id=image123", nil)
	rr := httptest.NewRecorder()

	handler.GetImage(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, "Failed to get image URL", resp.Error)
}

func TestGetImage_Success(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{
		getURLFunc: func(ctx context.Context, imageID string) (string, error) {
			return "http://cdn/" + imageID, nil
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/images?image_id=image123", nil)
	rr := httptest.NewRecorder()

	handler.GetImage(rr, req)

	require.Equal(t, http.StatusFound, rr.Code)
	require.Equal(t, "http://cdn/image123", rr.Header().Get("Location"))
}

func TestDeleteImage_Unauthorized(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := httptest.NewRequest(http.MethodDelete, "/images", nil)
	rr := httptest.NewRecorder()

	handler.DeleteImage(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteImage_MissingID(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := withUser(httptest.NewRequest(http.MethodDelete, "/images", nil), 2)
	rr := httptest.NewRecorder()

	handler.DeleteImage(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "image_id is required")
}

func TestDeleteImage_Error(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{
		deleteFunc: func(ctx context.Context, imageID string) error {
			return errors.New("boom")
		},
	})
	req := withUser(httptest.NewRequest(http.MethodDelete, "/images?image_id=image123", nil), 2)
	rr := httptest.NewRecorder()

	handler.DeleteImage(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	resp := decodeErrorResponse(t, rr.Body)
	require.Equal(t, "Failed to delete image", resp.Error)
}

func TestDeleteImage_Success(t *testing.T) {
	handler := NewHandler(&fakeImageUseCase{})
	req := withUser(httptest.NewRequest(http.MethodDelete, "/images?image_id=image123", nil), 2)
	rr := httptest.NewRecorder()

	handler.DeleteImage(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Empty(t, rr.Body.String())
}