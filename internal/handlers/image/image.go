package image

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	imageUC image.ImageUseCase
}

func NewHandler(imageUC image.ImageUseCase) *Handler {
	return &Handler{
		imageUC: imageUC,
	}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	return userID, ok
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputils.Error(w, r, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		httputils.Error(w, r, "Failed to get image from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		httputils.Error(w, r, "Invalid image format", http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/" + ext[1:]
	}

	imageID, err := h.imageUC.UploadImage(r.Context(), file, header.Filename, header.Size, contentType)
	if err != nil {
		httputils.InternalError(w, r, "Failed to upload image")
		return
	}

	url, err := h.imageUC.GetImageURL(r.Context(), imageID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get image URL")
		return
	}

	response := UploadImageResponse{
		ImageID: imageID,
		URL:     url,
	}

	httputils.Created(w, r, response)
}

func (h *Handler) GetImageURL(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	imageID := r.URL.Query().Get("image_id")
	if imageID == "" {
		httputils.Error(w, r, "image_id is required", http.StatusBadRequest)
		return
	}

	url, err := h.imageUC.GetImageURL(r.Context(), imageID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get image URL")
		return
	}

	response := ImageURLResponse{
		URL: url,
	}

	httputils.Success(w, r, response)
}

func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := r.URL.Query().Get("image_id")
	if imageID == "" {
		httputils.Error(w, r, "image_id is required", http.StatusBadRequest)
		return
	}

	url, err := h.imageUC.GetImageURL(r.Context(), imageID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get image URL")
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	imageID := r.URL.Query().Get("image_id")
	if imageID == "" {
		httputils.Error(w, r, "image_id is required", http.StatusBadRequest)
		return
	}

	err := h.imageUC.DeleteImage(r.Context(), imageID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to delete image")
		return
	}

	httputils.NoContent(w, r)
}
