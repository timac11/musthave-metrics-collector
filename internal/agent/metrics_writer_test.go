package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

func TestWrite(t *testing.T) {
	t.Run("should each metric", func(t *testing.T) {
		// Create test server to capture requests
		ctx := t.Context()
		requests := []string{}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Record the request URL
			requests = append(requests, r.URL.String())

			w.WriteHeader(http.StatusOK)
		}))

		// Test data
		value1 := 1.23
		value2 := 6.789
		metrics := []model.Metrics{
			{ID: "metric1", MType: model.Gauge, Value: &value1},
			{ID: "metric2", MType: model.Counter, Value: &value2},
		}

		url := strings.TrimPrefix(testServer.URL, "http://")
		defaultConfig := MetricsWriterConfig{URL: url, Attempts: 1, AttemptsInterval: 2, Mode: "http"}

		mw, err := newMetricsWriter(defaultConfig)

		require.NoError(t, err)

		mw.Write(ctx, metrics)

		// Verify requests were made
		require.Len(t, requests, 1)
		assert.Contains(t, requests[0], "/updates")

		testServer.Close()
	})
}

func TestWriteMetric(t *testing.T) {
	t.Run("should handle HTTP request failure", func(t *testing.T) {
		ctx := t.Context()
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
		}))

		// Test data
		value := 99.9
		metric := model.Metrics{
			ID:    "metric",
			MType: model.Gauge,
			Value: &value,
		}

		url := strings.TrimPrefix(testServer.URL, "http://")
		defaultConfig := MetricsWriterConfig{URL: url, Attempts: 1, AttemptsInterval: 2, Mode: "http"}

		require.NotPanics(t, func() {
			mw, _ := newMetricsWriter(defaultConfig)
			mw.writeMetric(ctx, metric)
		})

		testServer.Close()
	})

	t.Run("should format float values correctly in URL", func(t *testing.T) {
		ctx := t.Context()
		testCases := []struct {
			name     string
			value    float64
			expected string
		}{
			{"integer value", 42.0, "/update"},
			{"decimal value", 42.5, "/update"},
			{"small decimal", 0.123, "/update"},
			{"large value", 1234567.89, "/update"},
			{"negative value", -10.5, "/update"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test server
				var capturedURL string
				var metric model.Metrics
				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedURL = r.URL.String()
					json.NewDecoder(r.Body).Decode(&metric)
					w.WriteHeader(http.StatusOK)
				}))

				sendedMetric := model.Metrics{
					ID:    "test",
					MType: model.Gauge,
					Value: &tc.value,
				}

				url := strings.TrimPrefix(testServer.URL, "http://")
				defaultConfig := MetricsWriterConfig{URL: url, Attempts: 1, AttemptsInterval: 2, Mode: "http"}
				mw, err := newMetricsWriter(defaultConfig)

				require.NoError(t, err)

				mw.writeMetric(ctx, sendedMetric)

				assert.Equal(t, tc.expected, capturedURL)
				assert.Equal(t, sendedMetric.ID, metric.ID)
				assert.Equal(t, sendedMetric.Value, metric.Value)
				assert.Equal(t, sendedMetric.MType, metric.MType)

				testServer.Close()
			})
		}
	})
}
