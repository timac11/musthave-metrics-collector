package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

var memsMetrics = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
	"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
	"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
	"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
	"Sys", "TotalAlloc",
}

func TestCollectMemsMetrics(t *testing.T) {
	t.Run("should collect all runtime memory metrics", func(t *testing.T) {
		mc := newMetricsCollector()
		metrics := mc.collectRuntimeMemsMetrics()

		metricsMap := buildMetricMap(metrics)

		for _, expected := range memsMetrics {
			metric, exists := metricsMap[expected]
			assert.True(t, exists, "Metric %s should be present", expected)
			assert.Equal(t, model.Gauge, metric.MType, "Metric %s should have type Gauge", expected)
			assert.NotNil(t, metric.Value, "Metric %s should have a value", expected)
			assert.IsType(t, float64(0), *metric.Value, "Metric %s value should be float64", expected)
			assert.NotNil(t, metric.Value)
			assert.GreaterOrEqual(t, *metric.Value, float64(0),
				"Metric %s should not be negative", metric.ID)
		}

		// Verify we have the exact number of metrics
		assert.Len(t, metrics, len(memsMetrics))
	})
}

func TestCollectAdditionalMetrics(t *testing.T) {
	t.Run("should collect PollCount and RandomValue metrics", func(t *testing.T) {
		mc := newMetricsCollector()
		metrics, pollCount := mc.collectAdditionalMetrics()
		metricsMap := buildMetricMap(metrics)

		require.Len(t, metrics, 2, "Should return only 2 random metrics")

		pollCountMetric := metricsMap["PollCount"]
		require.NotNil(t, pollCountMetric, "PollCount metric should not be null")
		assert.Equal(t, model.Counter, pollCountMetric.MType)
		assert.NotNil(t, pollCountMetric.Delta)
		assert.Equal(t, int64(1), *pollCountMetric.Delta)
		assert.Equal(t, pollCount, *pollCountMetric.Delta)

		// Test RandomValue metric
		randomValueMetric := metricsMap["RandomValue"]
		require.NotNil(t, randomValueMetric, "RandomValue metric should not be null")
		assert.Equal(t, model.Gauge, randomValueMetric.MType)
		assert.NotNil(t, randomValueMetric.Value)
	})
}

func TestIntegration(t *testing.T) {
	t.Run("all metrics should have unique IDs", func(t *testing.T) {
		mc := newMetricsCollector()
		collectedValue := mc.Collect()

		ids := make(map[string]bool)
		for _, metric := range collectedValue.metrics {
			assert.False(t, ids[metric.ID], "Duplicate metric ID found: %s", metric.ID)
			ids[metric.ID] = true
		}
	})

	t.Run("metrics should have correct types", func(t *testing.T) {
		mc := newMetricsCollector()
		collectedValue := mc.Collect()

		for _, metric := range collectedValue.metrics {
			switch metric.ID {
			case "PollCount":
				assert.Equal(t, model.Counter, metric.MType,
					"Metric %s should be Counter type", metric.ID)
			default:
				assert.Equal(t, model.Gauge, metric.MType,
					"Metric %s should be Gauge type", metric.ID)
			}
		}
	})
}

func buildMetricMap(metrics []model.Metrics) map[string]model.Metrics {
	metricsMap := make(map[string]model.Metrics)
	for _, metric := range metrics {
		metricsMap[metric.ID] = metric
	}

	return metricsMap
}
