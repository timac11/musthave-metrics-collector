package agent

import (
	"math/rand"
	"runtime"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type MetricsCollector struct {
	pollCount int64
}

func (mc *MetricsCollector) Collect() []model.Metrics {
	memsMetrics := mc.collectMemsMetrics()
	additionalMetrics := mc.collectAdditionalMetrics()

	return append(memsMetrics, additionalMetrics...)
}

func (mc *MetricsCollector) collectMemsMetrics() []model.Metrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metricsMap := make(map[string]float64)

	metricsMap["Alloc"] = float64(stats.Alloc)
	metricsMap["BuckHashSys"] = float64(stats.BuckHashSys)
	metricsMap["Frees"] = float64(stats.Frees)
	metricsMap["GCCPUFraction"] = stats.GCCPUFraction
	metricsMap["GCSys"] = float64(stats.GCSys)
	metricsMap["HeapAlloc"] = float64(stats.HeapAlloc)
	metricsMap["HeapIdle"] = float64(stats.HeapIdle)
	metricsMap["HeapInuse"] = float64(stats.HeapInuse)
	metricsMap["HeapObjects"] = float64(stats.HeapObjects)
	metricsMap["HeapReleased"] = float64(stats.HeapReleased)
	metricsMap["HeapSys"] = float64(stats.HeapSys)
	metricsMap["LastGC"] = float64(stats.LastGC)
	metricsMap["Lookups"] = float64(stats.Lookups)
	metricsMap["MCacheInuse"] = float64(stats.MCacheInuse)
	metricsMap["MCacheSys"] = float64(stats.MCacheSys)
	metricsMap["MSpanInuse"] = float64(stats.MSpanInuse)
	metricsMap["MSpanSys"] = float64(stats.MSpanSys)
	metricsMap["Mallocs"] = float64(stats.Mallocs)
	metricsMap["NextGC"] = float64(stats.NextGC)
	metricsMap["NumForcedGC"] = float64(stats.NumForcedGC)
	metricsMap["NumGC"] = float64(stats.NumGC)
	metricsMap["OtherSys"] = float64(stats.OtherSys)
	metricsMap["PauseTotalNs"] = float64(stats.PauseTotalNs)
	metricsMap["StackInuse"] = float64(stats.StackInuse)
	metricsMap["StackSys"] = float64(stats.StackSys)
	metricsMap["Sys"] = float64(stats.Sys)
	metricsMap["TotalAlloc"] = float64(stats.TotalAlloc)

	var metrics []model.Metrics

	for key, value := range metricsMap {
		metrics = append(metrics, model.Metrics{ID: key, MType: model.Gauge, Value: &value})
	}

	return metrics
}

func (mc *MetricsCollector) collectAdditionalMetrics() []model.Metrics {
	mc.pollCount += 1
	pollCount := float64(mc.pollCount)
	randomValue := rand.Float64()

	metrics := []model.Metrics{
		{ID: "PollCount", MType: model.Counter, Value: &pollCount},
		{ID: "RandomValue", MType: model.Gauge, Value: &randomValue},
	}

	return metrics
}

func newMetricsCollector() *MetricsCollector {
	mc := &MetricsCollector{pollCount: 0}
	return mc
}
