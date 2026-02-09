package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

func TestPositiveServerSendScenario(t *testing.T) {
	writtenLogs := []*model.AuditLog{}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var log *model.AuditLog

		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&log)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writtenLogs = append(writtenLogs, log)

		w.WriteHeader(http.StatusOK)
	}))

	auditServerSender := NewAuditLogServerSender(testServer.URL)

	nowTime := time.Now().Unix()

	auditLog := model.AuditLog{
		TS:        nowTime,
		Metrics:   []string{"first"},
		IPAddress: "192.168.1.0",
	}

	err := auditServerSender.send(&auditLog)

	require.NoError(t, err)
	assert.Equal(t, len(writtenLogs), 1)
	assert.Equal(t, writtenLogs[0].TS, auditLog.TS)
	assert.Equal(t, writtenLogs[0].Metrics[0], auditLog.Metrics[0])
	assert.Equal(t, writtenLogs[0].IPAddress, auditLog.IPAddress)
}
