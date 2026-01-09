package main

import (
	"log"

	"github.com/timac11/musthave-metrics-collector/internal/agent"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func main() {
	flags := config.InitAgentConfig()
	logger.Initialize("INFO")

	metricsAgent, err := agent.NewMetricsAgent(flags)

	if err != nil {
		log.Fatal(err)
	}

	metricsAgent.Start()
}
