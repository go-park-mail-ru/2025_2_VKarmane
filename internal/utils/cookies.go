package utils

import (
	"net/http"
	"time"
)

func SetAuthCookie(w http.ResponseWriter, token string, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func ClearAuthCookie(w http.ResponseWriter, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}

func GetAuthCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetCSRFCookie(w http.ResponseWriter, token string, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
		Expires:  time.Now().Add(1 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func ClearCSRFCookie(w http.ResponseWriter, isProduction bool) {
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}

func GetCSRFCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
