package router

import (
	handler "github.com/timac11/musthave-metrics-collector/internal/handler"
	"net/http"
)

func InitRouter(mux *http.ServeMux) {
	mux.HandleFunc(`/update/`, updateMetricRoute)
	mux.HandleFunc(`/`, handler.NotFound) // default handler
}

func updateMetricRoute(res http.ResponseWriter, req *http.Request) {
	// valdiate request type
	if req.Method == http.MethodPost {
		handler.UpdateMetric(res, req)
	} else {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
