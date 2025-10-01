package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessAndCreated(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Success(rr, r, map[string]string{"ok": "1"})
	require.Equal(t, http.StatusOK, rr.Code)
	var m map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &m)
	assert.Equal(t, "1", m["ok"])

	rr2 := httptest.NewRecorder()
	Created(rr2, r, map[string]int{"id": 7})
	require.Equal(t, http.StatusCreated, rr2.Code)
}

func TestValidationErrors(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	verrs := utils.ValidationErrors{{Field: "login", Tag: "required", Value: "", Message: "Поле login обязательно для заполнения"}}
	ValidationErrors(rr, r, verrs)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestErrorHelpers(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/x", nil)

	t.Run("ValidationError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ValidationError(rr, r, "bad", "field")
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("UnauthorizedError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		UnauthorizedError(rr, r, "unauth", models.ErrCodeUnauthorized)
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("NotFoundError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		NotFoundError(rr, r, "nope")
		require.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("ConflictError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ConflictError(rr, r, "conflict", models.ErrCodeResourceConflict)
		require.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("InternalError", func(t *testing.T) {
		rr := httptest.NewRecorder()
		InternalError(rr, r, "boom")
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
