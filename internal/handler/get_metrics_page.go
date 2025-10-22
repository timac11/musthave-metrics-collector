package handler

import (
	"github.com/timac11/musthave-metrics-collector/internal/handler/templates"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	htmlTemplate "html/template"
	"net/http"
)

func (apiContainer *ApplicationAPIContainer) GetMetricsPage(res http.ResponseWriter, req *http.Request) {
	metrics := apiContainer.service.GetAll()

	indexTpl, err := htmlTemplate.ParseFS(templates.Index, "index.gohtml")

	logger.Debug("Template was parsed")

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = indexTpl.Execute(res, metrics)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Template was executed")
}
