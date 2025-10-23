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

			// Если origin пустой (например, от curl или прямых запросов), разрешаем
			// Это нужно для API тестирования и прямых запросов
			allowed := false
			// if origin == "" {
			// 	appLogger.Info("Empty origin, allowing for API testing", "origin", origin)
			// 	allowed = true
			// } else {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}
			// }

			// Устанавливаем CORS заголовки только для разрешенных origins
			if allowed {
				appLogger.Info("Origin allowed", "origin", origin)
				if origin == "" {
					// Для пустого origin (API тестирование) не устанавливаем credentials
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					// Для браузерных запросов устанавливаем точный origin и credentials
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			} else {
				appLogger.Warn("Origin not in allowed list, blocking request", "origin", origin)
				// Не устанавливаем CORS заголовки для неразрешенных origins
				// Это заблокирует запрос браузером
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			w.Header().Set("Access-Control-Max-Age", "3600")

			appLogger.Debug("CORS Headers set",
				"origin", w.Header().Get("Access-Control-Allow-Origin"),
				"methods", w.Header().Get("Access-Control-Allow-Methods"),
				"headers", w.Header().Get("Access-Control-Allow-Headers"))

			// Обрабатываем preflight запросы
			if r.Method == http.MethodOptions {
				appLogger.Info("Handling OPTIONS preflight request")
				w.WriteHeader(http.StatusOK)
				return
			}

			appLogger.Debug("Forwarding request to handler")
			next.ServeHTTP(w, r)
		})
	}
}
