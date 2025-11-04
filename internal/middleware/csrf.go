package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func CSRFMiddleware(authKey []byte) func(http.Handler) http.Handler {
	isProduction := os.Getenv("ENV") == "production"
	
	csrfProtect := csrf.Protect(
		authKey,
		csrf.Secure(isProduction),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Path("/"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "CSRF token validation failed", "code": "CSRF_TOKEN_INVALID"}`))
		})),
	)

	return csrfProtect
}

func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}
