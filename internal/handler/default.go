package handler

import (
	"net/http"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Not found", http.StatusNotFound)
}
