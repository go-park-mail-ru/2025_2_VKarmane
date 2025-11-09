package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_NoCookie(t *testing.T) {
	mw := AuthMiddleware("secret")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(next).ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mw := AuthMiddleware("secret")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "bad"})
	mw(next).ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ValidTokenSetsContext(t *testing.T) {
	mw := AuthMiddleware("secret")
	token, err := utils.GenerateJWT(12, "u", "secret")
	require.NoError(t, err)

	var gotUserID int
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := GetUserIDFromContext(r.Context())
		require.True(t, ok)
		gotUserID = id
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	mw(next).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, 12, gotUserID)
}
