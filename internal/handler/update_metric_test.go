package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
	"github.com/timac11/musthave-metrics-collector/internal/repository"
	"github.com/timac11/musthave-metrics-collector/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestPositiveUpdateMetricHandler(t *testing.T) {
	val := float64(1)

	dbFilePath := "./tmp/db.json"
	persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
	storage := repository.NewMemStorage(persistentStorage, false)
	service := service.NewService(storage)
	handlers := NewApplicationAPIContainer(*service)

	type result struct {
		code        int
		contentType string
	}

	tests := []struct {
		name   string
		metric model.Metrics
		result result
	}{
		{
			name: "positive test",
			metric: model.Metrics{
				ID:    "first",
				MType: model.Gauge,
				Value: &val,
			},
			result: result{
				code:        200,
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updateURL := "/update/" + test.metric.MType + "/" + test.metric.ID + "/" + strconv.FormatFloat(*test.metric.Value, 'f', 2, 64)
			request := httptest.NewRequest(http.MethodPost, updateURL, nil)
			w := httptest.NewRecorder()
			handlers.UpdateMetric(w, request)

			res := w.Result()
			assert.Equal(t, test.result.code, res.StatusCode)
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.result.contentType, res.Header.Get("Content-Type"))
		})
	}

	os.Remove(dbFilePath)
}

func TestNegativeUpdateMetric(t *testing.T) {
	type result struct {
		code int
	}
	dbFilePath := "/temp/db.json"
	defer os.Remove(dbFilePath)

	persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
	storage := repository.NewMemStorage(persistentStorage, false)
	service := service.NewService(storage)
	handlers := NewApplicationAPIContainer(*service)

	tests := []struct {
		name   string
		url    string
		result result
	}{
		{
			name: "invalid url",
			url:  "/update/invalid-url",
			result: result{
				code: http.StatusNotFound,
			},
		},
		{
			name: "invalid metric type test",
			url:  "/update/some-invalid-metric-type/name/1",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid gauge metric value test",
			url:  "/update/gauge/name/asd",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid counter metric value test",
			url:  "/update/counter/name/some-string",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid counter metric value test",
			url:  "/update/counter/name/some-string",
			result: result{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()
			handlers.UpdateMetric(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.result.code, res.StatusCode)
		})
	}

}
