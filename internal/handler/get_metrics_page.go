package handler

import (
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"net/http"
)

func (container *ApplicationAPIContainer) GetMetricsPage(res http.ResponseWriter, req *http.Request) {
	metrics, err := container.service.GetAll()

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		logger.Error("Internal server error", err.Error())
		return
	}

	indexTpl := container.templatesMap["index"]
	res.Header().Set("Content-Type", "text/html")

	if indexTpl == nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Error("Failed to load metrics page template")
		return
	}

	logger.Debug("metrics", metrics)
	err = indexTpl.Execute(res, metrics)

	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Error("Internal server error", err.Error())
		return
	}
	logger.Debug("Template was executed")
}
