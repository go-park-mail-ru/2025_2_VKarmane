package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	log := logger.NewSlogLogger()
	mw := RequestLoggerMiddleware(log)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestRequestLoggerMiddleware_POST(t *testing.T) {
	log := logger.NewSlogLogger()
	mw := RequestLoggerMiddleware(log)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusCreated, rr.Code)
}
