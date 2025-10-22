package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (container *ApplicationAPIContainer) GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	metric := container.service.Get(metricType, metricName)

	switch v := metric.(type) {
	case int64:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%d", v)
	case float64:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%f", v)
	default:
		res.WriteHeader(http.StatusInternalServerError)
	}
}
