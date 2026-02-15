package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

func TestPositiveAuditorCollectScenario(t *testing.T) {
	filePath := "/tmp/audit-writer.json"

	serverWrittenLogs := []*model.AuditLog{}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var log *model.AuditLog

		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&log)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serverWrittenLogs = append(serverWrittenLogs, log)

		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()
	defer os.Remove(filePath)

	auditor := NewAuditor(filePath, testServer.URL)

	// time to subscribe
	time.Sleep(100 * time.Millisecond)

	value1 := 99.9
	metric1 := model.Metrics{
		ID:    "metric1",
		MType: model.Gauge,
		Value: &value1,
	}

	value2 := int64(10)
	metric2 := model.Metrics{
		ID:    "metric2",
		MType: model.Gauge,
		Delta: &value2,
	}

	metrics := []*model.Metrics{&metric1, &metric2}
	ipAddress := "192.168.1.1"

	auditor.Collect(metrics, ipAddress)

	// time to complete writing
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, len(serverWrittenLogs))
	assert.NotNil(t, serverWrittenLogs[0].TS)
	assert.Equal(t, "metric1", serverWrittenLogs[0].Metrics[0])
	assert.Equal(t, "metric2", serverWrittenLogs[0].Metrics[1])
	assert.Equal(t, "192.168.1.1", serverWrittenLogs[0].IPAddress)
}
