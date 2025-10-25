package main

import (
	"sync"
	"time"

	config "github.com/timac11/musthave-metrics-collector/cmd/agent/config"
	agent "github.com/timac11/musthave-metrics-collector/internal/agent"
	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

const (
	ReadPeriod  = 2 * time.Second
	WritePeriod = 10 * time.Second
)

var (
	metrics []model.Metrics
	mu      sync.RWMutex
)

func collectMetrics(mc *agent.MetricsCollector) {
	for {
		time.Sleep(ReadPeriod)

		logger.Debug("Start collect metrics")

		mu.Lock()
		metrics = mc.Collect()
		mu.Unlock()

		logger.Debug("Complete collect metrics")
	}
}

func writeMetrics(mw *agent.MetricsWriter) {
	for {
		time.Sleep(WritePeriod)

		logger.Debug("Start write metrics")

		mu.RLock()
		mw.Write(metrics)
		mu.RUnlock()

		logger.Debug("Complete write metrics")
	}
}

func main() {
	flags := config.InitFlags()

	metricsCollector := agent.NewMetricsCollector()
	metricsWriter := agent.NewMetricsWriter(flags.Address)

	go collectMetrics(metricsCollector)
	go writeMetrics(metricsWriter)
	select {}
}
