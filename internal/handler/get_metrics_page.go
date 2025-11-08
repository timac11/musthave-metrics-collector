package handler

import (
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"net/http"
)

func (container *ApplicationAPIContainer) GetMetricsPage(res http.ResponseWriter, req *http.Request) {
	metrics := container.service.GetAll()
	indexTpl := container.templatesMap["index"]

	if indexTpl == nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	err := indexTpl.Execute(res, metrics)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		logger.Error("Internal server error", err.Error())
		return
	}

	logger.Debug("Template was executed")
}
