package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/config"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

type MetricsAgent struct {
	writer    *MetricsWriter
	collector *MetricsCollector
	ch        chan *CollectedMetrics
	config    *config.AgentConfig
}

// Start create two goroutines:
// to start metrics collection to start sending metrics to server
func (agent *MetricsAgent) Start() {
	go agent.collectMetrics()
	go agent.writeMetrics()
	select {}
}

func NewMetricsAgent(agentConfig *config.AgentConfig) (*MetricsAgent, error) {
	// count of sending metrics t server workers should be more or equal 1
	if agentConfig.RateLimit == 0 {
		return nil, fmt.Errorf("invalid count of workers: %d", agentConfig.RateLimit)
	}

	writerConfig := MetricsWriterConfig{
		Attempts:         agentConfig.RetryAttempts,
		AttemptsInterval: agentConfig.RetryInterval,
		SigningKey:       agentConfig.SigningKey,
		CryptoKey:        agentConfig.CryptoKey,
	}
	mc := newMetricsCollector()
	mw, err := newMetricsWriter(agentConfig.Address, writerConfig)
	if err != nil {
		return nil, err
	}

	ch := make(chan *CollectedMetrics, agentConfig.RateLimit)

	agent := &MetricsAgent{writer: mw, collector: mc, config: agentConfig, ch: ch}
	return agent, nil
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

	var wg sync.WaitGroup
	numWorkers := int(agent.config.RateLimit)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go func(worker int) {
			defer wg.Done()
			agent.runWriteMetricsWorker(worker)
		}(i)
	}

	wg.Wait()
}

func (agent *MetricsAgent) runWriteMetricsWorker(worker int) {
	writer := agent.writer
	collector := agent.collector

	for collectedMetrics := range agent.ch {
		logger.Debug("Start write metrics", "worker", worker)

		metrics := collectedMetrics.metrics
		pollCount := collectedMetrics.pollCount

		collector.addValueToPollCount(-pollCount)

		err := writer.Write(metrics)

		if err != nil {
			logger.Error("Failed write metrics", "worker", worker, err)
			collector.addValueToPollCount(pollCount)
		} else {
			logger.Debug("Complete write metrics", "worker", worker)
		}
	}
}
