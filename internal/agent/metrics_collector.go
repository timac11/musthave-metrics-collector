package agent

import (
	"math/rand"
	"runtime"
	"slices"
	"sync/atomic"
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/model"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CollectedMetrics is contract for metric collection
// metrics - is slice of collected runtime metrics
// pollCount - is the number of measurements taken before sending to server
// this value is equal to calculate PollCount metric
type CollectedMetrics struct {
	metrics   []model.Metrics
	pollCount int64
}

// MetricsCollector is structure of metric collector
// pollCount - is actual number of measurements taken before sending to server
type MetricsCollector struct {
	pollCount atomic.Int64
}

// Collect function get runtime memory metrics, PollCount and RandomValue metrics
func (mc *MetricsCollector) Collect() CollectedMetrics {
	memsMetrics := mc.collectRuntimeMemsMetrics()
	additionalMetrics, pollCount := mc.collectAdditionalMetrics()
	usageMetrics := mc.collectUsageMemsMetrics()

	metrics := slices.Concat(memsMetrics, additionalMetrics, usageMetrics)

	return CollectedMetrics{metrics: metrics, pollCount: pollCount}
}

// Reset pollCount value in collector
func (mc *MetricsCollector) Reset() {
	mc.pollCount.Store(0)
}

func (mc *MetricsCollector) collectRuntimeMemsMetrics() []model.Metrics {
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

func (mc *MetricsCollector) collectUsageMemsMetrics() []model.Metrics {
	metrics := []model.Metrics{}
	vmStat, err := mem.VirtualMemory()

	if err == nil {
		total := float64(vmStat.Total)
		free := float64(vmStat.Free)

		freeMetrics := []model.Metrics{
			{ID: "TotalMemory", Value: &total, MType: model.Gauge},
			{ID: "FreeMemory", Value: &free, MType: model.Gauge},
		}

		metrics = append(metrics, freeMetrics...)
	}

	percent, err := cpu.Percent(time.Millisecond, false)

	if err == nil {
		metrics = append(metrics, model.Metrics{
			ID: "CPUutilization1", Value: &percent[0], MType: model.Gauge,
		})
	}

	return metrics
}

func (mc *MetricsCollector) collectAdditionalMetrics() ([]model.Metrics, int64) {
	mc.pollCount.Add(1)
	pollCount := int64(mc.pollCount.Load())
	randomValue := rand.Float64()

	metrics := []model.Metrics{
		{ID: "PollCount", MType: model.Counter, Delta: &pollCount},
		{ID: "RandomValue", MType: model.Gauge, Value: &randomValue},
	}

	return metrics, pollCount
}

func (mc *MetricsCollector) addValueToPollCount(value int64) {
	mc.pollCount.Add(value)
}

func newMetricsCollector() *MetricsCollector {
	mc := &MetricsCollector{pollCount: atomic.Int64{}}
	return mc
}
