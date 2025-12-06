package handler

import (
	"net/http"
)

func (container *ApplicationAPIContainer) DbPing(res http.ResponseWriter, req *http.Request) {
	err := container.service.DbPing()
	res.Header().Set("Content-Type", "text/html")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
