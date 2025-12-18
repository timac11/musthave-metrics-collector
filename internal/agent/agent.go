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
	config    *config.AgentConfig
}

func (agent *MetricsAgent) Start() {
	go agent.collectMetrics()
	go agent.writeMetrics()
	select {}
}

func NewMetricsAgent(agentConfig *config.AgentConfig) *MetricsAgent {
	mc := newMetricsCollector()
	writerConfig := MetricsWriterConfig{Attempts: agentConfig.RetryAttempts, AttemptsInterval: agentConfig.RetryInterval}
	mw := newMetricsWriter(agentConfig.Address, writerConfig)
	mu := sync.Mutex{}
	agent := &MetricsAgent{writer: mw, collector: mc, config: agentConfig, mu: &mu}
	return agent
}

func (agent *MetricsAgent) collectMetrics() {
	logger.Info("Agent started collect metrics task")

	ticker := time.Tick(time.Duration(agent.config.PollInterval) * time.Second)

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
	logger.Info("Agent started write metrics task")
	ticker := time.Tick(time.Duration(agent.config.ReportInterval) * time.Second)

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
