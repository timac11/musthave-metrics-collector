package main

import (
	"log"

	"github.com/timac11/musthave-metrics-collector/internal/agent"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func main() {
	flags, err := config.InitAgentConfig()

	if err != nil {
		log.Fatal(err)
	}

	logger.Initialize("INFO")

	metricsAgent := agent.NewMetricsAgent(flags)
	metricsAgent.Start()
}
