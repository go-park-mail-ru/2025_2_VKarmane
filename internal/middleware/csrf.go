package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
)

var safeMethods = map[string]bool{
    http.MethodGet:     true,
    http.MethodHead:    true,
    http.MethodOptions: true,
    http.MethodTrace:   true,
}

func IsSafeMethod(method string) bool {
    return safeMethods[method]
}

func CSRFMiddleware(jwtSecret string) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !IsSafeMethod(r.Method){
				cookie, err := r.Cookie("csrf_token")
				if err != nil {
					http.Error(w, "CSRF token was not provided in cookies", http.StatusForbidden)
					return
				}
				headerToken := r.Header.Get("X-CSRF-Token")
				if (headerToken == "") {
					http.Error(w, "CSRF token was not provided in headers", http.StatusForbidden)
					return
				}
				cookieToken := cookie.Value

				if (cookieToken != headerToken) {
					http.Error(w, "Tokens do not match", http.StatusForbidden)
					return
				}
				
				 if _, err := utils.ValidateCSRF(headerToken, jwtSecret); err != nil {
					http.Error(w, "Invalid CSRF token", http.StatusForbidden)
					return
				 }
			} 
		ctx := context.WithValue(r.Context(), "csrf_secret", jwtSecret)
		next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}