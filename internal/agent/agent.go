package agent

import (
	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"time"
)

type MetricsAgent struct {
	writer    *MetricsWriter
	collector *MetricsCollector
	ch        chan *[]model.Metrics
	config    *config.AgentConfig
}

func (agent *MetricsAgent) Start() {
	go agent.collectMetrics()
	go agent.writeMetrics()
	select {}
}

func NewMetricsAgent(agentConfig *config.AgentConfig) *MetricsAgent {
	writerConfig := MetricsWriterConfig{Attempts: agentConfig.RetryAttempts, AttemptsInterval: agentConfig.RetryInterval, SigningKey: agentConfig.SigningKey}
	mc := newMetricsCollector()
	mw := newMetricsWriter(agentConfig.Address, writerConfig)

	ch := make(chan *[]model.Metrics, agentConfig.RateLimit)

	agent := &MetricsAgent{writer: mw, collector: mc, config: agentConfig, ch: ch}
	return agent
}

func (agent *MetricsAgent) collectMetrics() {
	logger.Info("Agent started collect metrics task")
	ticker := time.Tick(time.Duration(agent.config.PollInterval) * time.Second)
	collector := agent.collector

	for range ticker {
		logger.Debug("Start collect metrics")
		metrics := collector.Collect()
		agent.ch <- &metrics
		logger.Debug("Complete collect metrics")
	}
}

func (agent *MetricsAgent) writeMetrics() {
	logger.Info("Agent started write metrics task")
	writer := agent.writer
	collector := agent.collector

	for metrics := range agent.ch {
		logger.Debug("Start write metrics")
		writer.Write(*metrics)
		collector.Reset()
		logger.Debug("Complete write metrics")
	}
}
