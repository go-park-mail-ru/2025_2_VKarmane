package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddleware(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	log := logger.NewSlogLogger()
	mw := CORSMiddleware(allowedOrigins, log)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_OptionsRequest(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	log := logger.NewSlogLogger()
	mw := CORSMiddleware(allowedOrigins, log)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	require.False(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddleware_NotAllowedOrigin(t *testing.T) {
	allowedOrigins := []string{"http://localhost:3000"}
	log := logger.NewSlogLogger()
	mw := CORSMiddleware(allowedOrigins, log)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	rr := httptest.NewRecorder()

	mw(next).ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rr.Code)
}
