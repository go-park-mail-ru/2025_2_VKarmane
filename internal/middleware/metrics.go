package middleware

import (
	"net/http"
	"strconv"
	"time"

	metrics "github.com/go-park-mail-ru/2025_2_VKarmane/internal/metrics"
	"github.com/gorilla/mux"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &statusRecorder{ResponseWriter: w, status: 200}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		method := r.Method
		status := strconv.Itoa(rw.status)

		route := mux.CurrentRoute(r)
        pathTemplate := r.URL.Path

        if route != nil {
            if tpl, err := route.GetPathTemplate(); err == nil {
                pathTemplate = tpl
            }
        }

		metrics.HttpRequestsTotal.WithLabelValues(method, pathTemplate, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, pathTemplate, status).Observe(duration)
	})
}


type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
