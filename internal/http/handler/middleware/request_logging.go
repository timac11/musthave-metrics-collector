package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func (m *Middleware) RequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		recorder := &ResponseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}

		start := time.Now()
		h.ServeHTTP(recorder, r)
		end := time.Now()

		logger.Info(
			"Request info:",
			"method", r.Method,
			"uri", r.RequestURI,
			"duration", end.Sub(start),
		)

		logger.Info(
			"Response info",
			"status", recorder.status,
			"bytes", recorder.body.Len(),
		)
	})
}
