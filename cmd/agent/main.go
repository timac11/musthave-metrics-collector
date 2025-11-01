package main

import (
	"github.com/timac11/musthave-metrics-collector/internal/agent"
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func main() {
	flags := config.InitAgentFlags()
	logger.Log(flags.Collectnterval)
	logger.Log(flags.WriteInterval)
	metricsAgent := agent.NewMetricsAgent(flags)
	metricsAgent.Start()
}
