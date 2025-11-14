package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

func CORSMiddleware(allowedOrigins []string, appLogger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			appLogger.Info("CORS Request",
				"method", r.Method,
				"origin", origin,
				"path", r.URL.Path)
			appLogger.Debug("Allowed Origins", "origins", allowedOrigins)

			var allowed bool
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if r.Method == http.MethodOptions {
				// Для OPTIONS запросов разрешаем origin, если он в списке разрешенных
				if origin != "" {
					var originAllowed bool
					for _, allowedOrigin := range allowedOrigins {
						if origin == allowedOrigin {
							originAllowed = true
							break
						}
					}
					if originAllowed {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					} else {
						// Если origin не разрешен, но OPTIONS запрос - все равно отвечаем
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Access-Control-Allow-Credentials", "true")
						appLogger.Warn("OPTIONS request from non-allowed origin, but allowing for CORS preflight", "origin", origin)
					}
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}

				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
				w.Header().Set("Access-Control-Max-Age", "3600")

				appLogger.Info("Handling OPTIONS preflight request", "path", r.URL.Path, "origin", origin)
				w.WriteHeader(http.StatusOK)
				return
			}

			if allowed {
				appLogger.Info("Origin allowed", "origin", origin)
				if origin == "" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			} else {
				appLogger.Warn("Origin not in allowed list, blocking request", "origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
			w.Header().Set("Access-Control-Max-Age", "3600")

			appLogger.Debug("CORS Headers set",
				"origin", w.Header().Get("Access-Control-Allow-Origin"),
				"methods", w.Header().Get("Access-Control-Allow-Methods"),
				"headers", w.Header().Get("Access-Control-Allow-Headers"))

			appLogger.Debug("Forwarding request to handler")
			next.ServeHTTP(w, r)
		})
	}
}
