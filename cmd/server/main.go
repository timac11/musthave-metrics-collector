package main

import (
	"log"
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/router"
)

func main() {
	run()
}

func run() {
	conf := config.InitServerConfig()
	logger.Initialize("INFO")
	mux := router.InitRouter()
	err := http.ListenAndServe(conf.Address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
