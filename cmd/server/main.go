package main

import (
	"net/http"

	router "github.com/timac11/musthave-metrics-collector/internal/router"
)

func main() {
	run()
}

func run() {
	mux := router.InitRouter()
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
