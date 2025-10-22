package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

func TestWrite(t *testing.T) {
	t.Run("should each metric", func(t *testing.T) {
		// Create test server to capture requests
		requests := []string{}
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Record the request URL
			requests = append(requests, r.URL.String())
			w.WriteHeader(http.StatusOK)
		}))
		originalBaseURL := baseURL
		baseURL = testServer.URL

		// Test data
		value1 := 1.23
		value2 := 6.789
		metrics := []model.Metrics{
			{ID: "metric1", MType: model.Gauge, Value: &value1},
			{ID: "metric2", MType: model.Counter, Value: &value2},
		}

		// Call the function
		Write(metrics)

		// Verify requests were made
		require.Len(t, requests, 2)
		assert.Contains(t, requests[0], "/gauge/metric1/1.23")
		assert.Contains(t, requests[1], "/counter/metric2/6.789")

		testServer.Close()
		baseURL = originalBaseURL
	})
}

func TestWriteMetric(t *testing.T) {
	t.Run("should handle HTTP request failure", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
		}))

		originalBaseURL := baseURL
		baseURL = testServer.URL

		// Test data
		value := 99.9
		metric := model.Metrics{
			ID:    "metric",
			MType: model.Gauge,
			Value: &value,
		}

		require.NotPanics(t, func() {
			writeMetric(metric)
		})

		testServer.Close()
		baseURL = originalBaseURL
	})

	t.Run("should format float values correctly in URL", func(t *testing.T) {
		testCases := []struct {
			name     string
			value    float64
			expected string
		}{
			{"integer value", 42.0, "/update/gauge/test/42"},
			{"decimal value", 42.5, "/update/gauge/test/42.5"},
			{"small decimal", 0.123, "/update/gauge/test/0.123"},
			{"large value", 1234567.89, "/update/gauge/test/1234567.89"},
			{"negative value", -10.5, "/update/gauge/test/-10.5"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test server
				var capturedURL string
				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedURL = r.URL.String()
					w.WriteHeader(http.StatusOK)
				}))

				originalBaseURL := baseURL
				baseURL = testServer.URL

				metric := model.Metrics{
					ID:    "test",
					MType: model.Gauge,
					Value: &tc.value,
				}

				writeMetric(metric)
				assert.Equal(t, tc.expected, capturedURL)

				testServer.Close()
				baseURL = originalBaseURL
			})
		}
	})
}

func TestWriteNilMetricEdge(t *testing.T) {
	t.Run("should handle nil value", func(t *testing.T) {
		// Create test server
		requestMade := false
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestMade = true
			w.WriteHeader(http.StatusOK)
		}))

		originalBaseURL := baseURL
		baseURL = testServer.URL

		metric := model.Metrics{
			ID:    "nilValue",
			MType: model.Gauge,
			Value: nil,
		}

		require.Panics(t, func() {
			writeMetric(metric)
		})

		assert.False(t, requestMade)

		testServer.Close()
		baseURL = originalBaseURL
	})
}
