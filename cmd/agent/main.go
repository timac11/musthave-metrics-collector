package main

import (
	"sync"
	"time"

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

func readMetrics() {
	for {
		time.Sleep(ReadPeriod)

		logger.Debug("Start collect metrics")

		mu.Lock()
		metrics = agent.Collect()
		mu.Unlock()

		logger.Debug("Complete collect metrics")
	}
}

func writeMetrics() {
	for {
		time.Sleep(WritePeriod)

		logger.Debug("Start write metrics")

		mu.RLock()
		agent.Write(metrics)
		mu.RUnlock()

		logger.Debug("Complete write metrics")
	}
}

func main() {
	go readMetrics()
	go writeMetrics()
	select {}
}
