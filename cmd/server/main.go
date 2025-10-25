package main

import (
	"net/http"

	config "github.com/timac11/musthave-metrics-collector/cmd/server/config"
	router "github.com/timac11/musthave-metrics-collector/internal/router"
)

func main() {
	run()
}

func run() {
	flags := config.InitFlags()
	mux := router.InitRouter()
	err := http.ListenAndServe(flags.Address, mux)
	if err != nil {
		panic(err)
	}
}
