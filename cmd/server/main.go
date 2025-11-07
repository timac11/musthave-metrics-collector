package main

import (
	"log"
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/router"
)

func main() {
	run()
}

func run() {
	flags := config.InitServerFlags()
	mux := router.InitRouter()
	err := http.ListenAndServe(flags.Address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
