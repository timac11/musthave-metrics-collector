package handler

import (
	"net/http"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Not found", http.StatusNotFound)
}
