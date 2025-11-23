package middleware

import (
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"net/http"
	"time"
)

func (m *Middleware) RequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		h.ServeHTTP(ww, r)
		end := time.Now()

		logger.Info(
			"Request info:",
			"method", r.Method,
			"uri", r.RequestURI,
			"duration", end.Sub(start),
		)

		logger.Info(
			"Response info",
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
		)
	})
}
