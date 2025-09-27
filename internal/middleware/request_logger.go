package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

func RequestLoggerMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			clientIP := getClientIP(r)
			userID, isAuthenticated := GetUserIDFromContext(r.Context())
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			fields := []interface{}{
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", duration.String(),
				"ip", clientIP,
				"user_agent", r.UserAgent(),
			}

			if isAuthenticated {
				fields = append(fields, "user_id", userID)
			}

			if wrapped.size > 0 {
				fields = append(fields, "size", wrapped.size)
			}

			if wrapped.statusCode >= 400 {
				log.Error("HTTP Request", fields...)
			} else {
				log.Info("HTTP Request", fields...)
			}
		})
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Forwarded"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

func SecurityLoggerMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if strings.HasPrefix(r.URL.Path, "/api/v1/") && r.URL.Path != "/api/v1/auth/register" && r.URL.Path != "/api/v1/auth/login" {
				userID, isAuthenticated := GetUserIDFromContext(r.Context())

				if !isAuthenticated {
					log.Warn("Unauthorized access attempt",
						"ip", clientIP,
						"path", r.URL.Path,
						"method", r.Method,
						"user_agent", r.UserAgent(),
					)
				} else {
					log.Info("Authenticated request",
						"ip", clientIP,
						"user_id", userID,
						"path", r.URL.Path,
						"method", r.Method,
					)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
