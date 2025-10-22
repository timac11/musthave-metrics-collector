package handler

import (
	"net/http"
)

func (apiContainer *ApplicationAPIContainer) GetMetricsPage(res http.ResponseWriter, req *http.Request) {


	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}