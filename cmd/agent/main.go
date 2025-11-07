package main

import (
	"github.com/timac11/musthave-metrics-collector/internal/agent"
	"github.com/timac11/musthave-metrics-collector/internal/config"
)

func main() {
	flags := config.InitAgentConfig()
	metricsAgent := agent.NewMetricsAgent(flags)
	metricsAgent.Start()
}
