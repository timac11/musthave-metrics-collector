package handler

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func (apiContainer *ApplicationAPIContainer) GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	metric := apiContainer.service.Get(metricType, metricName)

	if metric != nil {
		resp, err := json.Marshal(metric)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
			res.Write(resp)
		}
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusNotFound)
	}
}