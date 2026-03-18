package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

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
	g, appCtx := errgroup.WithContext(context.Background())
	defer appCtx.Err()

	go agent.collectMetrics(appCtx)
	go agent.writeMetrics(appCtx)

	// graceful shutdown
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	g.Go(func() error {
		<-quitChan
		logger.Info("graceful shutdown signal")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err.Error())
	}
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

func (agent *MetricsAgent) collectMetrics(ctx context.Context) {
	logger.Info("Agent started collect metrics task")
	ticker := time.NewTicker(time.Duration(agent.config.PollInterval) * time.Second)
	collector := agent.collector

	for {
		select {
		case <-ticker.C:
			logger.Debug("Start collect metrics")
			metrics := collector.Collect()
			agent.ch <- &metrics
			logger.Debug("Complete collect metrics")
		case <-ctx.Done():
			return
		}
	}
}

func (agent *MetricsAgent) writeMetrics(ctx context.Context) {
	logger.Info("Agent started write metrics task")

	var wg sync.WaitGroup
	numWorkers := int(agent.config.RateLimit)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go func(worker int) {
			defer wg.Done()
			agent.runWriteMetricsWorker(ctx, worker)
		}(i)
	}

	wg.Wait()
}

func (agent *MetricsAgent) runWriteMetricsWorker(ctx context.Context, worker int) {
	writer := agent.writer
	collector := agent.collector

	for {
		select {
		case collectedMetrics := <-agent.ch:
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
		case <-ctx.Done():
			return
		}
	}
}
