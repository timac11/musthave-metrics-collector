package handler

import (
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func (container *ApplicationAPIContainer) DBPing(res http.ResponseWriter, req *http.Request) {
	err := container.service.DBPing()
	res.Header().Set("Content-Type", "text/html")
	if err != nil {
		logger.Error("Failed ping database", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
