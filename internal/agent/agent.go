package agent

import (
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"sync"
	"time"
)

type MetricsAgent struct {
	writer    *MetricsWriter
	collector *MetricsCollector
	mu        *sync.Mutex
	metrics   []model.Metrics
	config     *config.AgentConfig
}

func (agent *MetricsAgent) Start() {
	go agent.collectMetrics()
	go agent.writeMetrics()
	select {}
}

func NewMetricsAgent(flags *config.AgentConfig) *MetricsAgent {
	mc := newMetricsCollector()
	mw := newMetricsWriter(flags.Address)
	mu := sync.Mutex{}
	agent := &MetricsAgent{writer: mw, collector: mc, config: flags, mu: &mu}
	return agent
}

func (agent *MetricsAgent) collectMetrics() {
	ticker := time.Tick(time.Duration(agent.config.Collectnterval) * time.Second)

	for range ticker {
		logger.Debug("Start collect metrics")

		agent.mu.Lock()
		metrics := agent.collector.Collect()
		agent.metrics = metrics
		agent.mu.Unlock()

		logger.Debug("Complete collect metrics")
	}
}

func (agent *MetricsAgent) writeMetrics() {
	ticker := time.Tick(time.Duration(agent.config.WriteInterval) * time.Second)

	for range ticker {
		logger.Debug("Start write metrics")

		agent.mu.Lock()
		metrics := agent.metrics
		agent.writer.Write(metrics)
		agent.collector.Reset()
		agent.mu.Unlock()

		logger.Debug("Complete write metrics")
	}
}
