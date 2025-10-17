package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestHandler_Logout_Smoke(t *testing.T) {
	store, err := repository.NewStore()
	require.NoError(t, err)
	svc := service.NewService(store)
	uc := usecase.NewUseCase(svc, store, "secret")
	h := NewHandler(uc, logger.NewSlogLogger())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	h.authHandler.Logout(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
