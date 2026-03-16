package main

import (
	"log"
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/router"

	"github.com/timac11/musthave-metrics-collector/cmd/version"
)

func main() {
	run()
}

func run() {
	conf := config.InitServerConfig()
	logger.Initialize("INFO")
	mux, err := router.InitRouter(conf)

	if err != nil {
		log.Fatal(err)
	}

	version.Print()
	logger.Info("Starting server on address: ", conf.Address)

	err = http.ListenAndServe(conf.Address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
